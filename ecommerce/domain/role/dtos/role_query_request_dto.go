package dtos

import (
	"github.com/gin-gonic/gin"

	infradtos "go-service/infrastructure/dtos"
)

type RoleQueryRequestDto struct {
	infradtos.PaginationQueryRequestDto
}

func AssignRoleQueryRequestDto(c *gin.Context) *RoleQueryRequestDto {
	q := &RoleQueryRequestDto{}
	if err := c.ShouldBindQuery(q); err != nil {
		panic(gin.Error{Err: err, Type: gin.ErrorTypeBind})
	}
	return q
}
