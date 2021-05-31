package usecase

import (
	"context"
	"encoding/json"
	"github.com/MyFirstBabyTime/Server/domain"
	"github.com/MyFirstBabyTime/Server/tx"
	"github.com/pkg/errors"
	"net/http"
)

type expenditureUsecase struct {
	// expenditureRepository is repository interface about domain.ExpenditureRepository
	expenditureRepository domain.ExpenditureRepository

	// txHandler is used for handling transaction to begin & commit or rollback
	txHandler txHandler

	elasticSearch elasticSearch
}

func ExpenditureUsecase(
	er domain.ExpenditureRepository,
	th txHandler,
	es elasticSearch,
) *expenditureUsecase {
	return &expenditureUsecase{
		expenditureRepository: er,

		txHandler:     th,
		elasticSearch: es,
	}
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

type elasticSearch interface {
	Create(ctx context.Context, index string, s string) (err error)
}

func (eu *expenditureUsecase) ExpenditureRegistration(ctx context.Context, req *domain.Expenditure, babyUUIDs []string) (err error) {
	_tx, err := eu.txHandler.BeginTx(ctx, nil)
	if err != nil {
		err = errors.Wrap(err, "failed to begin transaction")
		return
	}

	switch err = eu.expenditureRepository.Store(_tx, req, babyUUIDs); tErr := err.(type) {
	case nil:
		break
	case domain.ErrInvalidModel:
		err = errors.Wrap(err, "Expenditure Store return invalid model")
		err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusInternalServerError}
		_ = eu.txHandler.Rollback(_tx)
	case domain.ErrNoReferencedRow:
		switch tErr.ForeignKey {
		case "parent_uuid":
			err = errors.New("this parent_uuid is not exist")
			err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusNotFound}
			_ = eu.txHandler.Rollback(_tx)
			return
		case "expenditure_uuid":
			err = errors.New("this expenditure_uuid is not exist")
			err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusNotFound}
			_ = eu.txHandler.Rollback(_tx)
			return
		case "baby_uuid":
			err = errors.New("this baby_uuid is not exist")
			err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusNotFound}
			_ = eu.txHandler.Rollback(_tx)
			return
		default:
			err = errors.New("expenditure_baby_tag return unexpected error")
			err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusInternalServerError}
			_ = eu.txHandler.Rollback(_tx)
			return
		}
	case domain.ErrEntryDuplicate:
		switch tErr.DuplicateKey {
		case "expenditure_baby_tag.PRIMARY":
			err = errors.New("expenditure_baby_tag Store duplicate value")
			err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusConflict}
			_ = eu.txHandler.Rollback(_tx)
			return
		default:
			err = errors.New("expenditure_baby_tag Store return unexpected duplicate error")
			err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusInternalServerError}
			_ = eu.txHandler.Rollback(_tx)
			return
		}
	default:
		err = errors.Wrap(err, "Expenditure Store return unexpected error")
		err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusInternalServerError}
		_ = eu.txHandler.Rollback(_tx)
		return
	}

	body, _ := esRequestBodyGenerator(req, babyUUIDs)
	err = eu.elasticSearch.Create(ctx, "Expenditure", body)

	if err != nil {
		err = errors.Wrap(err, "Expenditure Store return unexpected elasticSearch error")
		err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusInternalServerError}
		_ = eu.txHandler.Rollback(_tx)
		return
	}

	_ = eu.txHandler.Commit(_tx)
	return nil
}

func esRequestBodyGenerator(req *domain.Expenditure, babyUUIDS []string) (string, error) {
	data := make(map[string]interface{})

	data["UUID"] = req.UUID
	data["ParentUUID"] = req.ParentUUID
	data["Name"] = req.Name
	data["Amount"] = req.Amount
	data["Rating"] = req.Rating
	data["Link"] = req.Link
	data["baby"] = babyUUIDS

	body, err := json.Marshal(data)

	return string(body), err
}
