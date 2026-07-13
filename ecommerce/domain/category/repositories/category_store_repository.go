package repositories

import (
	"context"
	"log"

	"go-service/entities"
	categoryInterfaces "go-service/domain/category/interfaces"
	"go-service/infrastructure/exceptions"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryStoreRepository struct {
	db            *gorm.DB
	categoryModel *gorm.DB
}

func NewCategoryStoreRepository(db *gorm.DB) categoryInterfaces.CategoryStoreRepositoryInterface {
	return &CategoryStoreRepository{
		db:            db,
		categoryModel: db.Model(&entities.CategoryEntity{}),
	}
}

func (repo *CategoryStoreRepository) WithTransaction(tx *gorm.DB) categoryInterfaces.CategoryStoreRepositoryInterface {
	return &CategoryStoreRepository{
		db:            tx,
		categoryModel: tx.Model(&entities.CategoryEntity{}),
	}
}

func (repo *CategoryStoreRepository) Create(ctx context.Context, entity *entities.CategoryEntity) *entities.CategoryEntity {
	query := repo.categoryModel.WithContext(ctx)

	err := query.Create(entity).Error
	if err != nil {
		log.Println("Error create category:", err)
		panic(*exceptions.ServerErrorException(err))
	}

	return entity
}

func (repo *CategoryStoreRepository) Update(ctx context.Context, entity *entities.CategoryEntity) *entities.CategoryEntity {
	query := repo.categoryModel.WithContext(ctx)

	err := query.Where("id = ?", entity.Id).Save(entity).Error
	if err != nil {
		log.Println("Error update category:", err)
		panic(*exceptions.ServerErrorException(err))
	}

	return entity
}

func (repo *CategoryStoreRepository) Delete(ctx context.Context, id uuid.UUID) {
	query := repo.categoryModel.WithContext(ctx)

	err := query.Where("id = ?", id).Delete(&entities.CategoryEntity{}).Error
	if err != nil {
		log.Println("Error delete category:", err)
		panic(*exceptions.ServerErrorException(err))
	}
}
