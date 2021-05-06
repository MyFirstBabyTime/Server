package usecase

import (
	"context"
	"github.com/MyFirstBabyTime/Server/domain"
	"github.com/MyFirstBabyTime/Server/tx"
	"github.com/aws/aws-sdk-go/service/s3"
)

// childrenUsecase is used for usecase layer which implement domain.ChildrenUsecase interface
type childrenUsecase struct {
	domain.ChildrenUsecase
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
type childrenUsecaseConfig interface {}

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
