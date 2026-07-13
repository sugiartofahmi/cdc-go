package services

import (
	"context"

	"go-service/entities"
	categoryInterfaces "go-service/domain/category/interfaces"
	categoryConstants "go-service/domain/category/constants"
	categoryDtos "go-service/domain/category/dtos"
	"go-service/infrastructure/exceptions"
	"go-service/infrastructure/utils"
	infradtos "go-service/infrastructure/dtos"

	"github.com/google/uuid"
)

type CategoryService struct {
	categoryQueryRepository categoryInterfaces.CategoryQueryRepositoryInterface
	categoryStoreRepository categoryInterfaces.CategoryStoreRepositoryInterface
}

func NewCategoryService(
	categoryQueryRepository categoryInterfaces.CategoryQueryRepositoryInterface,
	categoryStoreRepository categoryInterfaces.CategoryStoreRepositoryInterface,
) categoryInterfaces.CategoryServiceInterface {
	return &CategoryService{
		categoryQueryRepository: categoryQueryRepository,
		categoryStoreRepository: categoryStoreRepository,
	}
}

func (s *CategoryService) Pagination(ctx context.Context, dto *categoryDtos.CategoryQueryRequestDto) *infradtos.PaginationResultDto[entities.CategoryEntity] {
	return s.categoryQueryRepository.Pagination(ctx, dto)
}

func (s *CategoryService) Detail(ctx context.Context, id uuid.UUID) *entities.CategoryEntity {
	category := s.categoryQueryRepository.FindOneById(ctx, id)
	if category == nil {
		panic(*exceptions.NotFoundException(categoryConstants.CATEGORY_NOT_FOUND))
	}
	return category
}

func (s *CategoryService) Create(ctx context.Context, dto *categoryDtos.CategoryCreateRequestDto) *entities.CategoryEntity {
	isNameAlreadyUsed := s.categoryQueryRepository.IsExistByName(ctx, dto.Name)
	if isNameAlreadyUsed {
		panic(*exceptions.ConflictException(categoryConstants.CATEGORY_NAME_ALREADY_EXISTS))
	}

	entity := &entities.CategoryEntity{
		Name: dto.Name,
		Slug: utils.GenerateSlug(dto.Name),
	}
	return s.categoryStoreRepository.Create(ctx, entity)
}

func (s *CategoryService) Update(ctx context.Context, id uuid.UUID, dto *categoryDtos.CategoryUpdateRequestDto) *entities.CategoryEntity {
	category := s.categoryQueryRepository.FindOneById(ctx, id)
	if category == nil {
		panic(*exceptions.NotFoundException(categoryConstants.CATEGORY_NOT_FOUND))
	}

	if dto.Name != nil {
		isNameAlreadyUsedByOther := s.categoryQueryRepository.IsExistByNameExcludeId(ctx, *dto.Name, id)
		if isNameAlreadyUsedByOther {
			panic(*exceptions.ConflictException(categoryConstants.CATEGORY_NAME_ALREADY_EXISTS))
		}
		category.Name = *dto.Name
		category.Slug = utils.GenerateSlug(*dto.Name)
	}

	return s.categoryStoreRepository.Update(ctx, category)
}

func (s *CategoryService) Delete(ctx context.Context, id uuid.UUID) {
	category := s.categoryQueryRepository.FindOneById(ctx, id)
	if category == nil {
		panic(*exceptions.NotFoundException(categoryConstants.CATEGORY_NOT_FOUND))
	}
	s.categoryStoreRepository.Delete(ctx, id)
}
