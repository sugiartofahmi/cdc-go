package interfaces

import (
	"context"

	"go-service/domain/product/dtos"
	infradtos "go-service/infrastructure/dtos"
)

type ProductQueryRepositoryInterface interface {
	Pagination(ctx context.Context, dto *dtos.ProductQueryRequestDto) *infradtos.PaginationResultDto[dtos.ProductResultDto]
}
