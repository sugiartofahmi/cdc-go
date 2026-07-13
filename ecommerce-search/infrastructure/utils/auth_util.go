package utils

import (
	"go-service/infrastructure/exceptions"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const AuthUserKey = "user"

func GetAuthUser(c *gin.Context) JWTUser {
	value, exists := c.Get(AuthUserKey)
	if !exists {
		panic(*exceptions.UnauthenticatedException("unauthenticated"))
	}

	user, ok := value.(JWTUser)
	if !ok {
		panic(*exceptions.UnauthenticatedException("invalid auth user in context"))
	}

	return user
}

func GetAuthUserId(c *gin.Context) uuid.UUID {
	return GetAuthUser(c).Id
}

func GetAuthRoleId(c *gin.Context) uuid.UUID {
	return GetAuthUser(c).RoleId
}

func GetAuthUserRoleName(c *gin.Context) string {
	return GetAuthUser(c).RoleName
}
