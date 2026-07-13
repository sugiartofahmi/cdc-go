package seeder

import (
	"encoding/json"
	"log"
	"os"

	"go-service/entities"
	"go-service/infrastructure/utils"

	"gorm.io/gorm"
)

type CategorySeeder struct{}

func NewCategorySeeder() *CategorySeeder {
	return &CategorySeeder{}
}

func (s *CategorySeeder) Handle(db *gorm.DB) error {
	data, err := os.ReadFile("seeder/files/categories.json")
	if err != nil {
		return err
	}

	var rows []struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(data, &rows); err != nil {
		return err
	}

	if err := db.Where("1 = 1").Delete(&entities.CategoryEntity{}).Error; err != nil {
		return err
	}

	for _, row := range rows {
		category := entities.CategoryEntity{
			Name: row.Name,
			Slug: utils.GenerateSlug(row.Name),
		}
		if err := db.Create(&category).Error; err != nil {
			return err
		}
	}

	log.Printf("CategorySeeder: inserted %d categories", len(rows))
	return nil
}
