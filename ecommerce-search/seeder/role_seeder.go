package seeder

import (
	"encoding/json"
	"log"
	"os"

	"go-service/entities"

	"gorm.io/gorm"
)

type RoleSeeder struct{}

func NewRoleSeeder() *RoleSeeder {
	return &RoleSeeder{}
}

func (s *RoleSeeder) Handle(db *gorm.DB) error {
	data, err := os.ReadFile("seeder/files/roles.json")
	if err != nil {
		return err
	}

	var rows []struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(data, &rows); err != nil {
		return err
	}

	if err := db.Where("1 = 1").Delete(&entities.RoleEntity{}).Error; err != nil {
		return err
	}

	for _, row := range rows {
		role := entities.RoleEntity{Name: row.Name}
		if err := db.Create(&role).Error; err != nil {
			return err
		}
	}

	log.Printf("RoleSeeder: inserted %d roles", len(rows))
	return nil
}
