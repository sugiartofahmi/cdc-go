package repositories

import (
	"context"
	"log"

	"go-service/entities"
	roleInterfaces "go-service/domain/role/interfaces"
	"go-service/infrastructure/exceptions"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleStoreRepository struct {
	db        *gorm.DB
	roleModel *gorm.DB
}

func NewRoleStoreRepository(db *gorm.DB) roleInterfaces.RoleStoreRepositoryInterface {
	return &RoleStoreRepository{
		db:        db,
		roleModel: db.Model(&entities.RoleEntity{}),
	}
}

func (repo *RoleStoreRepository) WithTransaction(tx *gorm.DB) roleInterfaces.RoleStoreRepositoryInterface {
	return &RoleStoreRepository{
		db:        tx,
		roleModel: tx.Model(&entities.RoleEntity{}),
	}
}

func (repo *RoleStoreRepository) Create(ctx context.Context, entity *entities.RoleEntity) *entities.RoleEntity {
	query := repo.roleModel.WithContext(ctx)

	err := query.Create(entity).Error
	if err != nil {
		log.Println("Error create role:", err)
		panic(*exceptions.ServerErrorException(err))
	}

	return entity
}

func (repo *RoleStoreRepository) Update(ctx context.Context, entity *entities.RoleEntity) *entities.RoleEntity {
	query := repo.roleModel.WithContext(ctx)

	err := query.Where("id = ?", entity.Id).Save(entity).Error
	if err != nil {
		log.Println("Error update role:", err)
		panic(*exceptions.ServerErrorException(err))
	}

	return entity
}

func (repo *RoleStoreRepository) Delete(ctx context.Context, id uuid.UUID) {
	query := repo.roleModel.WithContext(ctx)

	err := query.Where("id = ?", id).Delete(&entities.RoleEntity{}).Error
	if err != nil {
		log.Println("Error delete role:", err)
		panic(*exceptions.ServerErrorException(err))
	}
}
