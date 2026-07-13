package controller

import (
	"net/http"

	"go-service/domain/auth/constants"
	"go-service/domain/auth/dtos"
	authInterfaces "go-service/domain/auth/interfaces"
	"go-service/infrastructure/middlewares"
	"go-service/infrastructure/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService authInterfaces.AuthServiceInterface
}

func NewAuthController(router *gin.Engine, authService authInterfaces.AuthServiceInterface) {
	controller := &AuthController{authService: authService}

	authRoute := router.Group("/api/v1/auth")
	authRoute.POST("/login", middlewares.ValidateRequestJson[dtos.AuthLoginRequestDto](), controller.Login())
	authRoute.POST("/register", middlewares.ValidateRequestJson[dtos.AuthRegisterRequestDto](), controller.Register())
}

func (c *AuthController) Login() gin.HandlerFunc {
	return func(httpContext *gin.Context) {
		ctx := httpContext.Request.Context()
		dto := httpContext.MustGet(middlewares.RequestBodyJsonKey).(*dtos.AuthLoginRequestDto)
		result := c.authService.Login(ctx, dto)
		response := utils.SuccessResponse(http.StatusOK, constants.AUTH_LOGIN_SUCCESS, result)
		httpContext.JSON(http.StatusOK, response)
	}
}

func (c *AuthController) Register() gin.HandlerFunc {
	return func(httpContext *gin.Context) {
		ctx := httpContext.Request.Context()
		dto := httpContext.MustGet(middlewares.RequestBodyJsonKey).(*dtos.AuthRegisterRequestDto)
		result := c.authService.Register(ctx, dto)
		response := utils.SuccessResponse(http.StatusCreated, constants.AUTH_REGISTER_SUCCESS, result)
		httpContext.JSON(http.StatusCreated, response)
	}
}
