package interfaces

import (
	"context"

	"go-service/entities"
	categoryDtos "go-service/domain/category/dtos"
	infradtos "go-service/infrastructure/dtos"

	"github.com/google/uuid"
)

type CategoryServiceInterface interface {
	Pagination(ctx context.Context, dto *categoryDtos.CategoryQueryRequestDto) *infradtos.PaginationResultDto[entities.CategoryEntity]
	Detail(ctx context.Context, id uuid.UUID) *entities.CategoryEntity
	Create(ctx context.Context, dto *categoryDtos.CategoryCreateRequestDto) *entities.CategoryEntity
	Update(ctx context.Context, id uuid.UUID, dto *categoryDtos.CategoryUpdateRequestDto) *entities.CategoryEntity
	Delete(ctx context.Context, id uuid.UUID)
}
