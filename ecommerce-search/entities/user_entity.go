package entities

import (
	"time"

	"go-service/infrastructure/exceptions"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserEntity struct {
	Id        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"  json:"id"`
	Name      string         `gorm:"size:255;not null"                                json:"name"`
	Email     string         `gorm:"size:255;not null"                                json:"email"`
	Password  string         `gorm:"size:255;not null"                                json:"-"`
	Status    int            `gorm:"not null;default:1"                               json:"status"`
	RoleId    uuid.UUID      `gorm:"type:uuid;not null"                              json:"role_id"`
	Role      *RoleEntity    `gorm:"foreignKey:RoleId"                                json:"role"`
	CreatedAt *time.Time     `gorm:"autoCreateTime"                                   json:"created_at"`
	UpdatedAt *time.Time     `gorm:"autoUpdateTime"                                   json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"                                            json:"-"`
}

func (UserEntity) TableName() string { return "users" }

func (user *UserEntity) BeforeCreate(tx *gorm.DB) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(*exceptions.ServerErrorException(err))
	}
	user.Password = string(hash)
	return nil
}
