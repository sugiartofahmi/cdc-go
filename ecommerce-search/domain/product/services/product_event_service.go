package services

import (
	"context"
	"log"

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
	log.Printf("event handled: op=%s id=%s", dto.Operation, dto.Id)
	switch dto.Operation {
	case "c", "r":
		s.upsert(ctx, dto)
	case "u":
		isSoftDelete :=  dto.DeletedAt != nil
		if isSoftDelete {
			s.productStoreRepository.Delete(ctx, dto.Id)
		} else {
			s.upsert(ctx, dto)
		}
	case "d":
		s.productStoreRepository.Delete(ctx, dto.Id)
	}
}


func (s *ProductEventService) upsert(ctx context.Context, dto *productDtos.ProductEventDto) {
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
}