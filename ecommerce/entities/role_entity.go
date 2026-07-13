package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleEntity struct {
	Id        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"  json:"id"`
	Name      string         `gorm:"size:255;not null"                                json:"name"`
	CreatedAt *time.Time     `gorm:"autoCreateTime"                                   json:"created_at"`
	UpdatedAt *time.Time     `gorm:"autoUpdateTime"                                   json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"                                            json:"-"`
}

func (RoleEntity) TableName() string { return "roles" }
