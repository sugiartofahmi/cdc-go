package interfaces

import (
	"context"

	"go-service/entities"
)

type AuthUserStoreRepositoryInterface interface {
	Create(ctx context.Context, entity *entities.UserEntity) *entities.UserEntity
}
