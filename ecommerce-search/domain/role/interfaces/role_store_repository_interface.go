package interfaces

import (
	"context"

	"go-service/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleStoreRepositoryInterface interface {
	Create(ctx context.Context, entity *entities.RoleEntity) *entities.RoleEntity
	Update(ctx context.Context, entity *entities.RoleEntity) *entities.RoleEntity
	Delete(ctx context.Context, id uuid.UUID)
	WithTransaction(tx *gorm.DB) RoleStoreRepositoryInterface
}
