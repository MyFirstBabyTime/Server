package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// containerRedeployRequest is request for cloudMaintainerHandler.ContainerRedeploy
type containerRedeployRequest struct {
	CloudManagementKey string `json:"cloud_management_key" validate:"required"`
	Image              string `json:"image" validate:"required"`
}

func (r *containerRedeployRequest) BindFrom(c *gin.Context) error {
	return errors.Wrap(c.BindJSON(r), "falied to BindJSON")
}
