package interfaces

import (
	"context"

	"go-service/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryStoreRepositoryInterface interface {
	Create(ctx context.Context, entity *entities.CategoryEntity) *entities.CategoryEntity
	Update(ctx context.Context, entity *entities.CategoryEntity) *entities.CategoryEntity
	Delete(ctx context.Context, id uuid.UUID)
	WithTransaction(tx *gorm.DB) CategoryStoreRepositoryInterface
}
