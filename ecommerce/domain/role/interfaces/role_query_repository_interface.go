package interfaces

import (
	"context"

	"go-service/entities"
	roleDtos "go-service/domain/role/dtos"
	infradtos "go-service/infrastructure/dtos"

	"github.com/google/uuid"
)

type RoleQueryRepositoryInterface interface {
	Pagination(ctx context.Context, dto *roleDtos.RoleQueryRequestDto) *infradtos.PaginationResultDto[entities.RoleEntity]
	FindOneById(ctx context.Context, id uuid.UUID) *entities.RoleEntity
	FindOneByName(ctx context.Context, name string) *entities.RoleEntity
	IsExistByName(ctx context.Context, name string) bool
	IsExistByNameExcludeId(ctx context.Context, name string, excludeId uuid.UUID) bool
	IsExistById(ctx context.Context, id uuid.UUID) bool
}
