package dtos

import (
	"github.com/gin-gonic/gin"

	infradtos "go-service/infrastructure/dtos"
)

type UserQueryRequestDto struct {
	infradtos.PaginationQueryRequestDto
}

func AssignUserQueryRequestDto(c *gin.Context) *UserQueryRequestDto {
	q := &UserQueryRequestDto{}
	if err := c.ShouldBindQuery(q); err != nil {
		panic(gin.Error{Err: err, Type: gin.ErrorTypeBind})
	}
	return q
}
