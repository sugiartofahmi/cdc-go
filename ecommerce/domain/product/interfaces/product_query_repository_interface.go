package interfaces

import (
	"context"

	"go-service/entities"
	productDtos "go-service/domain/product/dtos"
	infradtos "go-service/infrastructure/dtos"

	"github.com/google/uuid"
)

type ProductQueryRepositoryInterface interface {
	Pagination(ctx context.Context, dto *productDtos.ProductQueryRequestDto) *infradtos.PaginationResultDto[entities.ProductEntity]
	FindOneById(ctx context.Context, id uuid.UUID) *entities.ProductEntity
	IsExistBySlug(ctx context.Context, slug string) bool
	IsExistBySlugExcludeId(ctx context.Context, slug string, excludeId uuid.UUID) bool
	IsExistById(ctx context.Context, id uuid.UUID) bool
}
