package dtos

import (
	"github.com/gin-gonic/gin"

	infradtos "go-service/infrastructure/dtos"
)

type CategoryQueryRequestDto struct {
	infradtos.PaginationQueryRequestDto
}

func AssignCategoryQueryRequestDto(c *gin.Context) *CategoryQueryRequestDto {
	q := &CategoryQueryRequestDto{}
	if err := c.ShouldBindQuery(q); err != nil {
		panic(gin.Error{Err: err, Type: gin.ErrorTypeBind})
	}
	return q
}
