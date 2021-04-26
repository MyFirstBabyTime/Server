package http

import (
	"github.com/MyFirstBabyTime/Server/domain"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

// cloudMaintainerHandler represent the http handler for article
type cloudMaintainerHandler struct {
	cUsecase  domain.CloudMaintainerUsecase
	validator validator
}

// validator is interface used for validating struct value
type validator interface {
	ValidateStruct(s interface{}) (err error)
}

// NewCloudMaintainerHandler will initialize the cloud-maintainer resources endpoint
func NewCloudMaintainerHandler(r *gin.Engine, cu domain.CloudMaintainerUsecase, v validator) {
	h := &cloudMaintainerHandler{
		cUsecase:  cu,
		validator: v,
	}

	r.POST("/redeploy", h.ContainerRedeploy)
}

// ContainerRedeploy deliver data to ContainerRedeploy of domain.CloudMaintainerUsecase
func (ch *cloudMaintainerHandler) ContainerRedeploy(c *gin.Context) {
	req := new(containerRedeployRequest)
	if err := ch.bindRequest(req, c); err != nil {
		c.JSON(http.StatusBadRequest, defaultResp(http.StatusBadRequest, 0, err.Error()))
		return
	}

	switch err := ch.cUsecase.ContainerRedeploy(c.Request.Context(), req.CloudManagementKey, req.Image); tErr := err.(type) {
	case nil:
		resp := defaultResp(http.StatusOK, 0, "succeed to container redeploy")
		c.JSON(http.StatusOK, resp)
	case domain.UsecaseError:
		c.JSON(tErr.Status, defaultResp(tErr.Status, tErr.Code, tErr.Error()))
	default:
		msg := errors.Wrap(err, "ContainerRedeploy return unexpected error").Error()
		c.JSON(http.StatusInternalServerError, defaultResp(http.StatusInternalServerError, 0, msg))
	}
	return
}

// bindRequest method bind *gin.Context to request having BindFrom method
func (ch *cloudMaintainerHandler) bindRequest(req interface {
	BindFrom(ctx *gin.Context) error
}, c *gin.Context) error {
	if err := req.BindFrom(c); err != nil {
		return errors.Wrap(err, "failed to bind req")
	}
	if err := ch.validator.ValidateStruct(req); err != nil {
		return errors.Wrap(err, "invalid request")
	}
	return nil
}

// defaultResp return response have status, code, message inform
func defaultResp(status, code int, msg string) (resp gin.H) {
	resp = gin.H{}
	resp["status"] = status
	resp["code"] = code
	resp["message"] = msg
	return
}