package http

import (
	"github.com/MyFirstBabyTime/Server/domain"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

//expenditureHandler represent the http handler for article
type expenditureHandler struct {
	eUsecase  domain.ExpenditureUsecase
	validator validator
	jwtHandler jwtHandler
}
