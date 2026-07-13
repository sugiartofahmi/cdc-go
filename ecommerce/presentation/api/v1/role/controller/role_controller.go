package controller

import (
	"net/http"

	"go-service/domain/role/constants"
	"go-service/domain/role/dtos"
	roleInterfaces "go-service/domain/role/interfaces"
	"go-service/infrastructure/middlewares"
	"go-service/infrastructure/utils"
	"go-service/infrastructure/validators"

	"github.com/gin-gonic/gin"
)

type RoleController struct {
	roleService roleInterfaces.RoleServiceInterface
}

func NewRoleController(router *gin.Engine, roleService roleInterfaces.RoleServiceInterface) {
	controller := &RoleController{roleService: roleService}

	roleRoute := router.Group("/api/v1/roles",
		middlewares.AuthorizationMiddleware(),
		middlewares.RoleMiddleware([]string{constants.ADMIN, constants.SUPER_ADMIN}),
	)
	roleRoute.GET("", controller.Pagination())
	roleRoute.GET("/:id", controller.Detail())
	roleRoute.POST("", middlewares.ValidateRequestJson[dtos.RoleCreateRequestDto](), controller.Create())
	roleRoute.PUT("/:id", middlewares.ValidateRequestJson[dtos.RoleUpdateRequestDto](), controller.Update())
	roleRoute.DELETE("/:id", controller.Delete())
}

func (c *RoleController) Pagination() gin.HandlerFunc {
	return func(httpContext *gin.Context) {
		ctx := httpContext.Request.Context()
		dto := dtos.AssignRoleQueryRequestDto(httpContext)
		result := c.roleService.Pagination(ctx, dto)
		meta := utils.PaginationMetaBuilder(dto.Page, dto.PerPage, int(result.Count))
		items := dtos.RoleResponseDtoFromEntities(result.Data)
		response := utils.SuccessResponsePagination(http.StatusOK, constants.ROLE_PAGINATION_SUCCESS, items, *meta)
		httpContext.JSON(http.StatusOK, response)
	}
}

func (c *RoleController) Detail() gin.HandlerFunc {
	return func(httpContext *gin.Context) {
		ctx := httpContext.Request.Context()
		id := validators.ValidateUUID(httpContext.Param("id"))
		result := dtos.RoleResponseDtoFromEntity(c.roleService.Detail(ctx, id))
		response := utils.SuccessResponse(http.StatusOK, constants.ROLE_DETAIL_SUCCESS, result)
		httpContext.JSON(http.StatusOK, response)
	}
}

func (c *RoleController) Create() gin.HandlerFunc {
	return func(httpContext *gin.Context) {
		ctx := httpContext.Request.Context()
		dto := httpContext.MustGet(middlewares.RequestBodyJsonKey).(*dtos.RoleCreateRequestDto)
		result := dtos.RoleResponseDtoFromEntity(c.roleService.Create(ctx, dto))
		response := utils.SuccessResponse(http.StatusCreated, constants.ROLE_CREATE_SUCCESS, result)
		httpContext.JSON(http.StatusCreated, response)
	}
}

func (c *RoleController) Update() gin.HandlerFunc {
	return func(httpContext *gin.Context) {
		ctx := httpContext.Request.Context()
		id := validators.ValidateUUID(httpContext.Param("id"))
		dto := httpContext.MustGet(middlewares.RequestBodyJsonKey).(*dtos.RoleUpdateRequestDto)
		result := dtos.RoleResponseDtoFromEntity(c.roleService.Update(ctx, id, dto))
		response := utils.SuccessResponse(http.StatusOK, constants.ROLE_UPDATE_SUCCESS, result)
		httpContext.JSON(http.StatusOK, response)
	}
}

func (c *RoleController) Delete() gin.HandlerFunc {
	return func(httpContext *gin.Context) {
		ctx := httpContext.Request.Context()
		id := validators.ValidateUUID(httpContext.Param("id"))
		c.roleService.Delete(ctx, id)
		response := utils.SuccessResponse(http.StatusOK, constants.ROLE_DELETE_SUCCESS, nil)
		httpContext.JSON(http.StatusOK, response)
	}
}
