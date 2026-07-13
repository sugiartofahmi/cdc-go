package interfaces

import (
	"context"

	"go-service/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserStoreRepositoryInterface interface {
	Create(ctx context.Context, entity *entities.UserEntity) *entities.UserEntity
	Update(ctx context.Context, entity *entities.UserEntity) *entities.UserEntity
	Delete(ctx context.Context, id uuid.UUID)
	WithTransaction(tx *gorm.DB) UserStoreRepositoryInterface
}
