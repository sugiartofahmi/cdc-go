package dtos

import (
	infradtos "go-service/infrastructure/dtos"

	"github.com/gin-gonic/gin"
)

type ProductQueryRequestDto struct {
	infradtos.PaginationQueryRequestDto

	Category string  `form:"category"`
	MinPrice float64 `form:"min_price"`
	MaxPrice float64 `form:"max_price"`
}

func AssignProductQueryRequestDto(c *gin.Context) *ProductQueryRequestDto {
	q := &ProductQueryRequestDto{}
	if err := c.ShouldBindQuery(q); err != nil {
		panic(gin.Error{
			Err:  err,
			Type: gin.ErrorTypeBind,
		})
	}
	return q
}
