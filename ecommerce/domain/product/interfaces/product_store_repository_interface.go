package interfaces

import (
	"context"

	"go-service/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductStoreRepositoryInterface interface {
	Create(ctx context.Context, entity *entities.ProductEntity) *entities.ProductEntity
	Update(ctx context.Context, entity *entities.ProductEntity) *entities.ProductEntity
	Delete(ctx context.Context, id uuid.UUID)
	WithTransaction(tx *gorm.DB) ProductStoreRepositoryInterface
}
