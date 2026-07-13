package repositories

import (
	"context"
	"log"

	"go-service/entities"
	productInterfaces "go-service/domain/product/interfaces"
	"go-service/infrastructure/exceptions"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductStoreRepository struct {
	db           *gorm.DB
	productModel *gorm.DB
}

func NewProductStoreRepository(db *gorm.DB) productInterfaces.ProductStoreRepositoryInterface {
	return &ProductStoreRepository{
		db:           db,
		productModel: db.Model(&entities.ProductEntity{}),
	}
}

func (repo *ProductStoreRepository) WithTransaction(tx *gorm.DB) productInterfaces.ProductStoreRepositoryInterface {
	return &ProductStoreRepository{
		db:           tx,
		productModel: tx.Model(&entities.ProductEntity{}),
	}
}

func (repo *ProductStoreRepository) Create(ctx context.Context, entity *entities.ProductEntity) *entities.ProductEntity {
	query := repo.productModel.WithContext(ctx)

	err := query.Create(entity).Error
	if err != nil {
		log.Println("Error create product:", err)
		panic(*exceptions.ServerErrorException(err))
	}

	return entity
}

func (repo *ProductStoreRepository) Update(ctx context.Context, entity *entities.ProductEntity) *entities.ProductEntity {
	query := repo.productModel.WithContext(ctx)

	err := query.Where("id = ?", entity.Id).Save(entity).Error
	if err != nil {
		log.Println("Error update product:", err)
		panic(*exceptions.ServerErrorException(err))
	}

	return entity
}

func (repo *ProductStoreRepository) Delete(ctx context.Context, id uuid.UUID) {
	query := repo.productModel.WithContext(ctx)

	err := query.Where("id = ?", id).Delete(&entities.ProductEntity{}).Error
	if err != nil {
		log.Println("Error delete product:", err)
		panic(*exceptions.ServerErrorException(err))
	}
}
