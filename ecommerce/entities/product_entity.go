package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductEntity struct {
	Id          uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name        string          `gorm:"size:255;not null"                                json:"name"`
	Slug        string          `gorm:"size:255;not null;uniqueIndex"                    json:"slug"`
	Description string          `gorm:"type:text"                                        json:"description"`
	Price       float64         `gorm:"type:numeric(15,2);not null"                      json:"price"`
	Stock       int             `gorm:"not null;default:0"                               json:"stock"`
	CategoryId  uuid.UUID       `gorm:"type:uuid;not null"                               json:"category_id"`
	Category    *CategoryEntity `gorm:"foreignKey:CategoryId"                            json:"category,omitempty"`
	CreatedAt   *time.Time      `gorm:"autoCreateTime"                                   json:"created_at"`
	UpdatedAt   *time.Time      `gorm:"autoUpdateTime"                                   json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index"                                            json:"-"`
}

func (ProductEntity) TableName() string { return "products" }
