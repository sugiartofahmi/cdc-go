package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryEntity struct {
	Id        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name      string         `gorm:"size:255;not null"                                json:"name"`
	Slug      string         `gorm:"size:255;not null;uniqueIndex"                    json:"slug"`
	CreatedAt *time.Time     `gorm:"autoCreateTime"                                   json:"created_at"`
	UpdatedAt *time.Time     `gorm:"autoUpdateTime"                                   json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"                                            json:"-"`
}

func (CategoryEntity) TableName() string { return "categories" }
