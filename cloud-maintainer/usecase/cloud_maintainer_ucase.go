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

// ContainerRedeploy is implement domain.CloudMaintainerUsecase interface
func (cu *cloudMaintainerUsecase) ContainerRedeploy(ctx context.Context, key string, image string) (err error) {
	if cu.cloudManagementKey != key {
		err = errors.New("cloudManagementKey is do not match")
		err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusForbidden}
		return
	}

	go func() {
		time.Sleep(time.Second)
		b, err := exec.Command("docker", "service", "update", "--image", image, "FirstBabyTime_server").Output()
		fmt.Println(string(b))
		fmt.Println(err)
	}()

	return nil
}