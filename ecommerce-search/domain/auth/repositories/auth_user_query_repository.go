package repositories

import (
	"context"
	"log"

	"go-service/entities"
	authInterfaces "go-service/domain/auth/interfaces"
	"go-service/infrastructure/exceptions"

	"gorm.io/gorm"
)

type AuthUserQueryRepository struct {
	db        *gorm.DB
	userModel *gorm.DB
}

func NewAuthUserQueryRepository(db *gorm.DB) authInterfaces.AuthUserQueryRepositoryInterface {
	return &AuthUserQueryRepository{
		db:        db,
		userModel: db.Model(&entities.UserEntity{}),
	}
}

func (r *AuthUserQueryRepository) FindOneByEmailWithRole(ctx context.Context, email string) *entities.UserEntity {
	var result entities.UserEntity
	err := r.userModel.WithContext(ctx).Preload("Role").Where("email = ?", email).First(&result).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	} else if err != nil {
		log.Println("Error find user by email with role:", err)
		panic(*exceptions.ServerErrorException(err))
	}
	return &result
}

func (r *AuthUserQueryRepository) IsExistsByEmail(ctx context.Context, email string) bool {
	var count int64
	err := r.userModel.WithContext(ctx).Where("email = ?", email).Count(&count).Error
	if err != nil {
		log.Println("Error check user exists by email:", err)
		panic(*exceptions.ServerErrorException(err))
	}
	return count > 0
}
