package utils

import (
	"gorm.io/gorm"

	infradtos "go-service/infrastructure/dtos"
)

func Paginate(q *infradtos.PaginationQueryRequestDto) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if q.Page <= 0 {
			q.Page = 1
		}

		if q.PerPage <= 0 || q.PerPage > 100 {
			q.PerPage = 10
		}

		offset := (q.Page - 1) * q.PerPage
		return db.Offset(offset).Limit(q.PerPage)
	}
}
