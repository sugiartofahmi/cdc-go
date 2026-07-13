package interfaces

import (
	"context"

	"go-service/entities"
	productDtos "go-service/domain/product/dtos"
	infradtos "go-service/infrastructure/dtos"

	"github.com/google/uuid"
)

type ProductServiceInterface interface {
	Pagination(ctx context.Context, dto *productDtos.ProductQueryRequestDto) *infradtos.PaginationResultDto[entities.ProductEntity]
	Detail(ctx context.Context, id uuid.UUID) *entities.ProductEntity
	Create(ctx context.Context, dto *productDtos.ProductCreateRequestDto) *entities.ProductEntity
	Update(ctx context.Context, id uuid.UUID, dto *productDtos.ProductUpdateRequestDto) *entities.ProductEntity
	Delete(ctx context.Context, id uuid.UUID)
}
