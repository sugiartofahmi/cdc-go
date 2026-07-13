package interfaces

import (
	"context"

	"go-service/entities"

	"github.com/google/uuid"
)

type AuthRoleQueryRepositoryInterface interface {
	IsExistById(ctx context.Context, id uuid.UUID) bool
	FindOneByName(ctx context.Context, name string) *entities.RoleEntity
}
