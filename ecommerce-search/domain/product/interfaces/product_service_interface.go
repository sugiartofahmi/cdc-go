package interfaces

import (
	"context"

	"go-service/domain/product/dtos"
	infraDtos "go-service/infrastructure/dtos"
)

type ProductServiceInterface interface {
	Pagination(ctx context.Context, dto *dtos.ProductQueryRequestDto) *infraDtos.PaginationResultDto[dtos.ProductResultDto]
}
