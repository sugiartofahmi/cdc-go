package repositories

import (
	"context"
	"log"

	authInterfaces "go-service/domain/auth/interfaces"
	"go-service/entities"
	"go-service/infrastructure/exceptions"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthRoleQueryRepository struct {
	db        *gorm.DB
	roleModel *gorm.DB
}

func NewAuthRoleQueryRepository(db *gorm.DB) authInterfaces.AuthRoleQueryRepositoryInterface {
	return &AuthRoleQueryRepository{
		db:        db,
		roleModel: db.Model(&entities.RoleEntity{}),
	}
}

func (r *AuthRoleQueryRepository) IsExistById(ctx context.Context, id uuid.UUID) bool {
	var count int64
	err := r.roleModel.WithContext(ctx).Where("id = ?", id).Count(&count).Error
	if err != nil {
		log.Println("Error check role exists by id:", err)
		panic(*exceptions.ServerErrorException(err))
	}
	return count > 0
}

func (r *AuthRoleQueryRepository) FindOneByName(ctx context.Context, name string) *entities.RoleEntity {
	var role entities.RoleEntity
	err := r.roleModel.WithContext(ctx).Where("name = ?", name).First(&role).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		log.Println("Error find role by name:", err)
		panic(*exceptions.ServerErrorException(err))
	}
	return &role
}
