package seeder

import (
	"encoding/json"
	"log"
	"os"

	"go-service/entities"
	"go-service/infrastructure/utils"

	"gorm.io/gorm"
)

type ProductSeeder struct{}

func NewProductSeeder() *ProductSeeder {
	return &ProductSeeder{}
}

func (s *ProductSeeder) Handle(db *gorm.DB) error {
	data, err := os.ReadFile("seeder/files/products.json")
	if err != nil {
		return err
	}

	var rows []struct {
		Name         string  `json:"name"`
		Description  string  `json:"description"`
		Price        float64 `json:"price"`
		Stock        int     `json:"stock"`
		CategoryName string  `json:"category_name"`
	}
	if err := json.Unmarshal(data, &rows); err != nil {
		return err
	}

	if err := db.Where("1 = 1").Delete(&entities.ProductEntity{}).Error; err != nil {
		return err
	}

	for _, row := range rows {
		var category entities.CategoryEntity
		if err := db.Where("name = ?", row.CategoryName).First(&category).Error; err != nil {
			log.Printf("ProductSeeder: skip product '%s', category '%s' not found", row.Name, row.CategoryName)
			continue
		}

		product := entities.ProductEntity{
			Name:        row.Name,
			Slug:        utils.GenerateSlug(row.Name),
			Description: row.Description,
			Price:       row.Price,
			Stock:       row.Stock,
			CategoryId:  category.Id,
		}
		if err := db.Create(&product).Error; err != nil {
			return err
		}
	}

	log.Printf("ProductSeeder: inserted products")
	return nil
}
