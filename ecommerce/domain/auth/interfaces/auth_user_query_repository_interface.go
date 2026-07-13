package interfaces

import (
	"context"

	"go-service/entities"
)

type AuthUserQueryRepositoryInterface interface {
	FindOneByEmailWithRole(ctx context.Context, email string) *entities.UserEntity
	IsExistsByEmail(ctx context.Context, email string) bool
}
