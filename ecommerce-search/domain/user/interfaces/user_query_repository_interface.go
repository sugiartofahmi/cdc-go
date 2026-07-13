package interfaces

import (
	"context"

	"go-service/entities"
	userDtos "go-service/domain/user/dtos"
	infradtos "go-service/infrastructure/dtos"

	"github.com/google/uuid"
)

type UserQueryRepositoryInterface interface {
	Pagination(ctx context.Context, dto *userDtos.UserQueryRequestDto) *infradtos.PaginationResultDto[entities.UserEntity]
	FindOneById(ctx context.Context, id uuid.UUID) *entities.UserEntity
	FindOneByIdForDetail(ctx context.Context, id uuid.UUID) *entities.UserEntity
	IsExistByName(ctx context.Context, name string) bool
	IsExistByNameExcludeId(ctx context.Context, name string, excludeId uuid.UUID) bool
	IsExistByEmail(ctx context.Context, email string) bool
	IsExistById(ctx context.Context, id uuid.UUID) bool
}
