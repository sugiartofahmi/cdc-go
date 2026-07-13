package utils

import (
	infradtos "go-service/infrastructure/dtos"
	"go-service/infrastructure/enums"
)

func Paginate(q *infradtos.PaginationQueryRequestDto) {
	if q.Page <= 0 {
		q.Page = 1
	}

	if q.PerPage <= 0 || q.PerPage > 100 {
		q.PerPage = 10
	}

	if q.SortBy == "" {
		q.SortBy = "created_at"
	}

	if q.Order == "" {
		q.Order = enums.SortOrderDesc
	}
}
