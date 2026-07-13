package interfaces

import (
	"context"

	"github.com/google/uuid"
)

type UserRoleQueryRepositoryInterface interface {
	IsExistById(ctx context.Context, id uuid.UUID) bool
}
