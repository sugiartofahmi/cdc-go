package singleton

import "gorm.io/gorm"

func PostgresSingleton() *gorm.DB {
	return db
}
