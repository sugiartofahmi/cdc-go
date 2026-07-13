package middlewares

import (
	"github.com/gin-gonic/gin"
)

const RequestBodyJsonKey = "RequestBodyJsonKey"

func ValidateRequestJson[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var dto T
		if err := c.ShouldBindJSON(&dto); err != nil {
			panic(gin.Error{
				Err:  err,
				Type: gin.ErrorTypeBind,
			})
		}

		c.Set(RequestBodyJsonKey, &dto)
		c.Next()
	}
}
