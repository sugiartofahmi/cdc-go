package dtos

import (
	"go-service/infrastructure/enums"

	"github.com/gin-gonic/gin"
)

type PaginationQueryRequestDto struct {
	Search  string              `form:"search"`
	PerPage int                 `form:"per_page"`
	Page    int                 `form:"page"`
	SortBy  string              `form:"sort_by"`
	Order   enums.SortOrderEnum `form:"order"`
}

func AssignPaginationQueryRequestDto(c *gin.Context) *PaginationQueryRequestDto {
	q := &PaginationQueryRequestDto{}
	if err := c.ShouldBindQuery(q); err != nil {
		panic(gin.Error{
			Err:  err,
			Type: gin.ErrorTypeBind,
		})
	}
	return q
}
