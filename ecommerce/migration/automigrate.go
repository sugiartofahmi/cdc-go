package migration

import (
	"gorm.io/gorm"

	"go-service/entities"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entities.CategoryEntity{},
		&entities.ProductEntity{},
	)
}
