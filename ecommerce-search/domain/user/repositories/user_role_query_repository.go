package repositories

import (
	"context"
	"log"

	userInterfaces "go-service/domain/user/interfaces"
	"go-service/entities"
	"go-service/infrastructure/exceptions"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRoleQueryRepository struct {
	db        *gorm.DB
	roleModel *gorm.DB
}

func NewUserRoleQueryRepository(db *gorm.DB) userInterfaces.UserRoleQueryRepositoryInterface {
	return &UserRoleQueryRepository{
		db:        db,
		roleModel: db.Model(&entities.RoleEntity{}),
	}
}

func (r *UserRoleQueryRepository) IsExistById(ctx context.Context, id uuid.UUID) bool {
	var count int64
	err := r.roleModel.WithContext(ctx).Where("id = ?", id).Count(&count).Error
	if err != nil {
		log.Println("Error check role exists by id:", err)
		panic(*exceptions.ServerErrorException(err))
	}
	return count > 0
}
