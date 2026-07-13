package seeder

import "gorm.io/gorm"

type Seeder interface {
	Handle(db *gorm.DB) error
}
