package interfaces

import (
	"context"

	"go-service/domain/product/dtos"
)

type ProductStoreRepositoryInterface interface {
	Upsert(ctx context.Context, dto *dtos.ProductResultDto)
	Delete(ctx context.Context, id string)
}
