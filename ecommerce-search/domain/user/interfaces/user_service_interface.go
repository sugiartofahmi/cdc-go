package interfaces

import (
	"context"

	"go-service/entities"
	userDtos "go-service/domain/user/dtos"
	infradtos "go-service/infrastructure/dtos"

	"github.com/google/uuid"
)

type UserServiceInterface interface {
	Pagination(ctx context.Context, dto *userDtos.UserQueryRequestDto) *infradtos.PaginationResultDto[entities.UserEntity]
	Detail(ctx context.Context, id uuid.UUID) *entities.UserEntity
	Create(ctx context.Context, dto *userDtos.UserCreateRequestDto) *entities.UserEntity
	Update(ctx context.Context, id uuid.UUID, dto *userDtos.UserUpdateRequestDto) *entities.UserEntity
	Delete(ctx context.Context, id uuid.UUID)
}
