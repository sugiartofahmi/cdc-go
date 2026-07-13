package controller

import (
	"net/http"

	"go-service/domain/role/constants"
	"go-service/domain/user/dtos"
	userConstants "go-service/domain/user/constants"
	userInterfaces "go-service/domain/user/interfaces"
	"go-service/infrastructure/middlewares"
	"go-service/infrastructure/utils"
	"go-service/infrastructure/validators"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService userInterfaces.UserServiceInterface
}

func NewUserController(router *gin.Engine, userService userInterfaces.UserServiceInterface) {
	controller := &UserController{userService: userService}

	userRoute := router.Group("/api/v1/users",
		middlewares.AuthorizationMiddleware(),
		middlewares.RoleMiddleware([]string{constants.ADMIN, constants.SUPER_ADMIN}),
	)
	userRoute.GET("", controller.Pagination())
	userRoute.GET("/:id", controller.Detail())
	userRoute.POST("", middlewares.ValidateRequestJson[dtos.UserCreateRequestDto](), controller.Create())
	userRoute.PUT("/:id", middlewares.ValidateRequestJson[dtos.UserUpdateRequestDto](), controller.Update())
	userRoute.DELETE("/:id", controller.Delete())
}

func (c *UserController) Pagination() gin.HandlerFunc {
	return func(httpContext *gin.Context) {
		ctx := httpContext.Request.Context()
		dto := dtos.AssignUserQueryRequestDto(httpContext)
		result := c.userService.Pagination(ctx, dto)
		meta := utils.PaginationMetaBuilder(dto.Page, dto.PerPage, int(result.Count))
		items := dtos.UserResponseDtoFromEntities(result.Data)
		response := utils.SuccessResponsePagination(http.StatusOK, userConstants.USER_PAGINATION_SUCCESS, items, *meta)
		httpContext.JSON(http.StatusOK, response)
	}
}

func (c *UserController) Detail() gin.HandlerFunc {
	return func(httpContext *gin.Context) {
		ctx := httpContext.Request.Context()
		id := validators.ValidateUUID(httpContext.Param("id"))
		result := dtos.UserResponseDtoFromEntity(c.userService.Detail(ctx, id))
		response := utils.SuccessResponse(http.StatusOK, userConstants.USER_DETAIL_SUCCESS, result)
		httpContext.JSON(http.StatusOK, response)
	}
}

func (c *UserController) Create() gin.HandlerFunc {
	return func(httpContext *gin.Context) {
		ctx := httpContext.Request.Context()
		dto := httpContext.MustGet(middlewares.RequestBodyJsonKey).(*dtos.UserCreateRequestDto)
		result := dtos.UserResponseDtoFromEntity(c.userService.Create(ctx, dto))
		response := utils.SuccessResponse(http.StatusCreated, userConstants.USER_CREATE_SUCCESS, result)
		httpContext.JSON(http.StatusCreated, response)
	}
}

func (c *UserController) Update() gin.HandlerFunc {
	return func(httpContext *gin.Context) {
		ctx := httpContext.Request.Context()
		id := validators.ValidateUUID(httpContext.Param("id"))
		dto := httpContext.MustGet(middlewares.RequestBodyJsonKey).(*dtos.UserUpdateRequestDto)
		result := dtos.UserResponseDtoFromEntity(c.userService.Update(ctx, id, dto))
		response := utils.SuccessResponse(http.StatusOK, userConstants.USER_UPDATE_SUCCESS, result)
		httpContext.JSON(http.StatusOK, response)
	}
}

func (c *UserController) Delete() gin.HandlerFunc {
	return func(httpContext *gin.Context) {
		ctx := httpContext.Request.Context()
		id := validators.ValidateUUID(httpContext.Param("id"))
		c.userService.Delete(ctx, id)
		response := utils.SuccessResponse(http.StatusOK, userConstants.USER_DELETE_SUCCESS, nil)
		httpContext.JSON(http.StatusOK, response)
	}
}
