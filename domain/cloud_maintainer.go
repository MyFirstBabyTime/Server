package domain

import (
	"context"
)

// CloudMaintainerUsecase is abstract interface about usecase layer using in delivery layer
type CloudMaintainerUsecase interface {
	ContainerRedeploy(ctx context.Context, key string, image string) (err error)
}