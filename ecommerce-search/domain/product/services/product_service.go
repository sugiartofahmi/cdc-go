package services

import (
	"context"

	"go-service/domain/product/dtos"
	productInterfaces "go-service/domain/product/interfaces"
	infradtos "go-service/infrastructure/dtos"
)

type ProductService struct {
	productQueryRepository productInterfaces.ProductQueryRepositoryInterface
}

func NewProductService(productQueryRepository productInterfaces.ProductQueryRepositoryInterface) productInterfaces.ProductServiceInterface {
	return &ProductService{
		productQueryRepository: productQueryRepository,
	}
}

func (s *ProductService) Pagination(ctx context.Context, dto *dtos.ProductQueryRequestDto) *infradtos.PaginationResultDto[dtos.ProductResultDto] {
	return s.productQueryRepository.Pagination(ctx, dto)
}
