package usecase

import (
	"context"
	"fmt"
	"github.com/MyFirstBabyTime/Server/domain"
	"github.com/pkg/errors"
	"net/http"
	"os/exec"
	"time"
)

// cloudMaintainerUsecase is used for usecase layer which implement domain.CloudMaintainerUsecase interface
type cloudMaintainerUsecase struct {
	// cloudManagementKey is used for check valid user before redeploy
	cloudManagementKey string
}

// CloudMaintainerUsecase return implementation of domain.CloudMaintainerUsecase
func CloudMaintainerUsecase(cloudManagementKey string) domain.CloudMaintainerUsecase {
	return &cloudMaintainerUsecase{
		cloudManagementKey: cloudManagementKey,
	}
}
