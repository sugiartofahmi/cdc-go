package interfaces

import (
	"context"

	"go-service/entities"
	categoryDtos "go-service/domain/category/dtos"
	infradtos "go-service/infrastructure/dtos"

	"github.com/google/uuid"
)

type CategoryQueryRepositoryInterface interface {
	Pagination(ctx context.Context, dto *categoryDtos.CategoryQueryRequestDto) *infradtos.PaginationResultDto[entities.CategoryEntity]
	FindOneById(ctx context.Context, id uuid.UUID) *entities.CategoryEntity
	FindOneBySlug(ctx context.Context, slug string) *entities.CategoryEntity
	IsExistByName(ctx context.Context, name string) bool
	IsExistByNameExcludeId(ctx context.Context, name string, excludeId uuid.UUID) bool
	IsExistById(ctx context.Context, id uuid.UUID) bool
}
