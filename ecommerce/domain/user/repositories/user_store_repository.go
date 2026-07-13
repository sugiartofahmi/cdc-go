package repositories

import (
	"context"
	"log"

	"go-service/entities"
	userInterfaces "go-service/domain/user/interfaces"
	"go-service/infrastructure/exceptions"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserStoreRepository struct {
	db        *gorm.DB
	userModel *gorm.DB
}

func NewUserStoreRepository(db *gorm.DB) userInterfaces.UserStoreRepositoryInterface {
	return &UserStoreRepository{
		db:        db,
		userModel: db.Model(&entities.UserEntity{}),
	}
}

func (repo *UserStoreRepository) WithTransaction(tx *gorm.DB) userInterfaces.UserStoreRepositoryInterface {
	return &UserStoreRepository{
		db:        tx,
		userModel: tx.Model(&entities.UserEntity{}),
	}
}

func (repo *UserStoreRepository) Create(ctx context.Context, entity *entities.UserEntity) *entities.UserEntity {
	query := repo.userModel.WithContext(ctx)

	err := query.Create(entity).Error
	if err != nil {
		log.Println("Error create user:", err)
		panic(*exceptions.ServerErrorException(err))
	}

	return entity
}

func (repo *UserStoreRepository) Update(ctx context.Context, entity *entities.UserEntity) *entities.UserEntity {
	query := repo.userModel.WithContext(ctx)

	err := query.Where("id = ?", entity.Id).Save(entity).Error
	if err != nil {
		log.Println("Error update user:", err)
		panic(*exceptions.ServerErrorException(err))
	}

	return entity
}

func (repo *UserStoreRepository) Delete(ctx context.Context, id uuid.UUID) {
	query := repo.userModel.WithContext(ctx)

	err := query.Where("id = ?", id).Delete(&entities.UserEntity{}).Error
	if err != nil {
		log.Println("Error delete user:", err)
		panic(*exceptions.ServerErrorException(err))
	}
}
