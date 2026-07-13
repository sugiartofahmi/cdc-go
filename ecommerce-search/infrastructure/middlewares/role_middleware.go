package middlewares

import (
	"slices"

	"github.com/gin-gonic/gin"

	"go-service/infrastructure/exceptions"
	"go-service/infrastructure/utils"
)

func RoleMiddleware(roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := utils.GetAuthUser(c)
		isRoleExists := slices.Contains(roles, user.RoleName)
		if !isRoleExists {
			panic(*exceptions.ForbiddenException("access denied"))
		}
		c.Next()
	}
}
