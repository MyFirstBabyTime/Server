package usecase

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/pkg/errors"
	"mime/multipart"
	"net/http"

	"github.com/MyFirstBabyTime/Server/domain"
	"github.com/MyFirstBabyTime/Server/tx"
	"github.com/aws/aws-sdk-go/service/s3"
)

// childrenUsecase is used for usecase layer which implement domain.ChildrenUsecase interface
type childrenUsecase struct {
	// myCfg is used for get config value for children usecase
	myCfg childrenUsecaseConfig

	// childrenRepository is repository interface about domain.Children model
	childrenRepository domain.ChildrenRepository

	// txHandler is used for handling transaction to begin & commit or rollback
	txHandler txHandler

	// s3Agency is used as agency about aws s3 API
	s3Agency s3Agency
}

// ChildrenUsecase return implementation of domain.ChildrenUsecase
func ChildrenUsecase(
	cfg childrenUsecaseConfig,
	cr domain.ChildrenRepository,
	th txHandler,
	sa s3Agency,
) domain.ChildrenUsecase {
	return &childrenUsecase{
		myCfg: cfg,

		childrenRepository: cr,

		txHandler: th,
		s3Agency:  sa,
	}
}

// childrenUsecaseConfig is interface get config value for children usecase
type childrenUsecaseConfig interface {
	// ChildrenProfileS3Bucket return aws s3 bucket name for children profile
	ChildrenProfileS3Bucket() string
}

// txHandler is used for handling transaction to begin & commit or rollback
type txHandler interface {
	// BeginTx method start transaction (get option from ctx)
	BeginTx(ctx context.Context, opts interface{}) (tx tx.Context, err error)

	// Commit method commit transaction
	Commit(tx tx.Context) (err error)

	// Rollback method rollback transaction
	Rollback(tx tx.Context) (err error)
}

// s3Agency is agency that agent various API about aws s3
type s3Agency interface {
	// PutObject method put(insert or update) object to s3
	PutObject(input *s3.PutObjectInput) (output *s3.PutObjectOutput, err error)
}

func (cu *childrenUsecase) CreateNewChildren(ctx context.Context, c *domain.Children, profile *multipart.FileHeader) (uuid string, err error) {
	_tx, err := cu.txHandler.BeginTx(ctx, nil)
	if err != nil {
		err = errors.Wrap(err, "failed to begin transaction")
		return
	}

	if domain.StringValue(c.UUID) == "" {
		if c.UUID, err = cu.childrenRepository.GetAvailableUUID(_tx); err != nil {
			err = domain.UsecaseError{UsecaseErr: errors.Wrap(err, "failed to GetAvailableUUID"), Status: http.StatusInternalServerError}
			_ = cu.txHandler.Rollback(_tx)
			return
		}
	}

	if profile != nil {
		c.ProfileUri = domain.String(c.GenerateProfileUri())
	}

	switch err = cu.childrenRepository.Store(_tx, c); tErr := err.(type) {
	case nil:
		break
	case domain.ErrInvalidModel:
		err = errors.Wrap(err, "children Store return invalid model")
		err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusInternalServerError}
		_ = cu.txHandler.Rollback(_tx)
		return
	case domain.ErrNoReferencedRow:
		switch tErr.ForeignKey {
		case "parent_uuid":
			err = errors.New("parent with that uuid is not exist")
			err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusNotFound}
			_ = cu.txHandler.Rollback(_tx)
			return
		default:
			err = errors.Wrap(err, "children Store return unexpected no referenced error")
			err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusInternalServerError}
			_ = cu.txHandler.Rollback(_tx)
			return
		}
	default:
		err = errors.Wrap(err, "children Store return unexpected error")
		err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusInternalServerError}
		_ = cu.txHandler.Rollback(_tx)
		return
	}

	if profile != nil {
		b := make([]byte, profile.Size)
		file, _ := profile.Open()
		defer func() { _ = file.Close() }()
		_, _ = file.Read(b)

		if _, err = cu.s3Agency.PutObject(&s3.PutObjectInput{
			Bucket: aws.String(cu.myCfg.ChildrenProfileS3Bucket()),
			Key:    aws.String(c.GenerateProfileUri()),
			Body:   bytes.NewReader(b),
			ACL:    aws.String("public-read"),
		}); err != nil {
			err = errors.Wrap(err, "s3 PutObject return unexpected error")
			err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusInternalServerError}
			_ = cu.txHandler.Rollback(_tx)
			return
		}
	}

	uuid = domain.StringValue(c.UUID)
	err = nil
	_ = cu.txHandler.Commit(_tx)
	return
}
