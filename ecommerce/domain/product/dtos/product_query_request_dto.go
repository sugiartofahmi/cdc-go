package dtos

import (
	"github.com/gin-gonic/gin"

	infradtos "go-service/infrastructure/dtos"
)

type ProductQueryRequestDto struct {
	infradtos.PaginationQueryRequestDto
}

func AssignProductQueryRequestDto(c *gin.Context) *ProductQueryRequestDto {
	q := &ProductQueryRequestDto{}
	if err := c.ShouldBindQuery(q); err != nil {
		panic(gin.Error{Err: err, Type: gin.ErrorTypeBind})
	}
	return q
}
