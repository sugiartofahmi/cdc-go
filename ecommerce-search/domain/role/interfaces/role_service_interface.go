package interfaces

import (
	"context"

	"go-service/entities"
	roleDtos "go-service/domain/role/dtos"
	infradtos "go-service/infrastructure/dtos"

	"github.com/google/uuid"
)

type RoleServiceInterface interface {
	Pagination(ctx context.Context, dto *roleDtos.RoleQueryRequestDto) *infradtos.PaginationResultDto[entities.RoleEntity]
	Detail(ctx context.Context, id uuid.UUID) *entities.RoleEntity
	Create(ctx context.Context, dto *roleDtos.RoleCreateRequestDto) *entities.RoleEntity
	Update(ctx context.Context, id uuid.UUID, dto *roleDtos.RoleUpdateRequestDto) *entities.RoleEntity
	Delete(ctx context.Context, id uuid.UUID)
}
