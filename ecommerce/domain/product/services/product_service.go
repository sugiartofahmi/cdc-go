package services

import (
	"context"

	"go-service/entities"
	categoryInterfaces "go-service/domain/category/interfaces"
	productInterfaces "go-service/domain/product/interfaces"
	productConstants "go-service/domain/product/constants"
	productDtos "go-service/domain/product/dtos"
	"go-service/infrastructure/exceptions"
	"go-service/infrastructure/utils"
	infradtos "go-service/infrastructure/dtos"

	"github.com/google/uuid"
)

type ProductService struct {
	productQueryRepository  productInterfaces.ProductQueryRepositoryInterface
	productStoreRepository  productInterfaces.ProductStoreRepositoryInterface
	categoryQueryRepository categoryInterfaces.CategoryQueryRepositoryInterface
}

func NewProductService(
	productQueryRepository productInterfaces.ProductQueryRepositoryInterface,
	productStoreRepository productInterfaces.ProductStoreRepositoryInterface,
	categoryQueryRepository categoryInterfaces.CategoryQueryRepositoryInterface,
) productInterfaces.ProductServiceInterface {
	return &ProductService{
		productQueryRepository:  productQueryRepository,
		productStoreRepository:  productStoreRepository,
		categoryQueryRepository: categoryQueryRepository,
	}
}

func (s *ProductService) Pagination(ctx context.Context, dto *productDtos.ProductQueryRequestDto) *infradtos.PaginationResultDto[entities.ProductEntity] {
	return s.productQueryRepository.Pagination(ctx, dto)
}

func (s *ProductService) Detail(ctx context.Context, id uuid.UUID) *entities.ProductEntity {
	product := s.productQueryRepository.FindOneById(ctx, id)
	if product == nil {
		panic(*exceptions.NotFoundException(productConstants.PRODUCT_NOT_FOUND))
	}
	return product
}

func (s *ProductService) Create(ctx context.Context, dto *productDtos.ProductCreateRequestDto) *entities.ProductEntity {
	if !s.categoryQueryRepository.IsExistById(ctx, dto.CategoryId) {
		panic(*exceptions.NotFoundException(productConstants.CATEGORY_NOT_FOUND))
	}

	slug := utils.GenerateSlug(dto.Name)
	if s.productQueryRepository.IsExistBySlug(ctx, slug) {
		slug = slug + "-" + uuid.New().String()[:8]
	}

	entity := &entities.ProductEntity{
		Name:        dto.Name,
		Slug:        slug,
		Description: dto.Description,
		Price:       dto.Price,
		Stock:       dto.Stock,
		CategoryId:  dto.CategoryId,
	}
	return s.productStoreRepository.Create(ctx, entity)
}

func (s *ProductService) Update(ctx context.Context, id uuid.UUID, dto *productDtos.ProductUpdateRequestDto) *entities.ProductEntity {
	product := s.productQueryRepository.FindOneById(ctx, id)
	if product == nil {
		panic(*exceptions.NotFoundException(productConstants.PRODUCT_NOT_FOUND))
	}

	if dto.Name != nil {
		slug := utils.GenerateSlug(*dto.Name)
		if s.productQueryRepository.IsExistBySlugExcludeId(ctx, slug, id) {
			slug = slug + "-" + uuid.New().String()[:8]
		}
		product.Name = *dto.Name
		product.Slug = slug
	}

	if dto.Description != nil {
		product.Description = *dto.Description
	}

	if dto.Price != nil {
		product.Price = *dto.Price
	}

	if dto.Stock != nil {
		product.Stock = *dto.Stock
	}

	if dto.CategoryId != nil {
		if !s.categoryQueryRepository.IsExistById(ctx, *dto.CategoryId) {
			panic(*exceptions.NotFoundException(productConstants.CATEGORY_NOT_FOUND))
		}
		product.CategoryId = *dto.CategoryId
	}

	return s.productStoreRepository.Update(ctx, product)
}

func (s *ProductService) Delete(ctx context.Context, id uuid.UUID) {
	product := s.productQueryRepository.FindOneById(ctx, id)
	if product == nil {
		panic(*exceptions.NotFoundException(productConstants.PRODUCT_NOT_FOUND))
	}
	s.productStoreRepository.Delete(ctx, id)
}
