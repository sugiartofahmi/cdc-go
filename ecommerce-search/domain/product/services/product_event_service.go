package services

import (
	"context"

	productDtos "go-service/domain/product/dtos"
	productInterfaces "go-service/domain/product/interfaces"
)

type ProductEventService struct {
	productStoreRepository productInterfaces.ProductStoreRepositoryInterface
}

func NewProductEventService(productStoreRepository productInterfaces.ProductStoreRepositoryInterface) productInterfaces.ProductEventServiceInterface {
	return &ProductEventService{
		productStoreRepository: productStoreRepository,
	}
}

func (s *ProductEventService) Handle(ctx context.Context, dto *productDtos.ProductEventDto) {
	switch dto.Operation {
	case "c", "r", "u":
		s.productStoreRepository.Upsert(ctx, &productDtos.ProductResultDto{
			Id:           dto.Id,
			Name:         dto.Name,
			Slug:         dto.Slug,
			Description:  dto.Description,
			Price:        dto.Price,
			Stock:        dto.Stock,
			CategoryId:   dto.CategoryId,
			CategoryName: "",
			CategorySlug: "",
			CreatedAt:    dto.CreatedAt,
			UpdatedAt:    dto.UpdatedAt,
		})
	case "d":
		s.productStoreRepository.Delete(ctx, dto.Id)
	}
}
