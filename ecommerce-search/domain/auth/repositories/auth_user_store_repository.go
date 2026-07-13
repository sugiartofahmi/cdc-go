package repositories

import (
	"context"
	"log"

	"go-service/entities"
	authInterfaces "go-service/domain/auth/interfaces"
	"go-service/infrastructure/exceptions"

	"gorm.io/gorm"
)

type AuthUserStoreRepository struct {
	db        *gorm.DB
	userModel *gorm.DB
}

func NewAuthUserStoreRepository(db *gorm.DB) authInterfaces.AuthUserStoreRepositoryInterface {
	return &AuthUserStoreRepository{
		db:        db,
		userModel: db.Model(&entities.UserEntity{}),
	}
}

func (r *AuthUserStoreRepository) Create(ctx context.Context, entity *entities.UserEntity) *entities.UserEntity {
	query := r.userModel.WithContext(ctx)

	err := query.Create(entity).Error
	if err != nil {
		log.Println("Error create user:", err)
		panic(*exceptions.ServerErrorException(err))
	}

	return entity
}
