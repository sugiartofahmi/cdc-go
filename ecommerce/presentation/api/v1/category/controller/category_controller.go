package controller

import (
	"net/http"

	categoryConstants "go-service/domain/category/constants"
	categoryDtos "go-service/domain/category/dtos"
	categoryInterfaces "go-service/domain/category/interfaces"
	"go-service/infrastructure/middlewares"
	"go-service/infrastructure/utils"
	"go-service/infrastructure/validators"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	categoryService categoryInterfaces.CategoryServiceInterface
}

func NewCategoryController(router *gin.Engine, categoryService categoryInterfaces.CategoryServiceInterface) {
	controller := &CategoryController{categoryService: categoryService}

	categoryRoute := router.Group("/api/v1/categories")
	categoryRoute.GET("", controller.Pagination())
	categoryRoute.GET("/:id", controller.Detail())
	categoryRoute.POST("", middlewares.ValidateRequestJson[categoryDtos.CategoryCreateRequestDto](), controller.Create())
	categoryRoute.PUT("/:id", middlewares.ValidateRequestJson[categoryDtos.CategoryUpdateRequestDto](), controller.Update())
	categoryRoute.DELETE("/:id", controller.Delete())
}

func (c *CategoryController) Pagination() gin.HandlerFunc {
	return func(httpContext *gin.Context) {
		ctx := httpContext.Request.Context()
		dto := categoryDtos.AssignCategoryQueryRequestDto(httpContext)
		result := c.categoryService.Pagination(ctx, dto)
		meta := utils.PaginationMetaBuilder(dto.Page, dto.PerPage, int(result.Count))
		items := categoryDtos.CategoryResponseDtoFromEntities(result.Data)
		response := utils.SuccessResponsePagination(http.StatusOK, categoryConstants.CATEGORY_PAGINATION_SUCCESS, items, *meta)
		httpContext.JSON(http.StatusOK, response)
	}
}

func (c *CategoryController) Detail() gin.HandlerFunc {
	return func(httpContext *gin.Context) {
		ctx := httpContext.Request.Context()
		id := validators.ValidateUUID(httpContext.Param("id"))
		result := categoryDtos.CategoryResponseDtoFromEntity(c.categoryService.Detail(ctx, id))
		response := utils.SuccessResponse(http.StatusOK, categoryConstants.CATEGORY_DETAIL_SUCCESS, result)
		httpContext.JSON(http.StatusOK, response)
	}
}

func (c *CategoryController) Create() gin.HandlerFunc {
	return func(httpContext *gin.Context) {
		ctx := httpContext.Request.Context()
		dto := httpContext.MustGet(middlewares.RequestBodyJsonKey).(*categoryDtos.CategoryCreateRequestDto)
		result := categoryDtos.CategoryResponseDtoFromEntity(c.categoryService.Create(ctx, dto))
		response := utils.SuccessResponse(http.StatusCreated, categoryConstants.CATEGORY_CREATE_SUCCESS, result)
		httpContext.JSON(http.StatusCreated, response)
	}
}

func (c *CategoryController) Update() gin.HandlerFunc {
	return func(httpContext *gin.Context) {
		ctx := httpContext.Request.Context()
		id := validators.ValidateUUID(httpContext.Param("id"))
		dto := httpContext.MustGet(middlewares.RequestBodyJsonKey).(*categoryDtos.CategoryUpdateRequestDto)
		result := categoryDtos.CategoryResponseDtoFromEntity(c.categoryService.Update(ctx, id, dto))
		response := utils.SuccessResponse(http.StatusOK, categoryConstants.CATEGORY_UPDATE_SUCCESS, result)
		httpContext.JSON(http.StatusOK, response)
	}
}

func (c *CategoryController) Delete() gin.HandlerFunc {
	return func(httpContext *gin.Context) {
		ctx := httpContext.Request.Context()
		id := validators.ValidateUUID(httpContext.Param("id"))
		c.categoryService.Delete(ctx, id)
		response := utils.SuccessResponse(http.StatusOK, categoryConstants.CATEGORY_DELETE_SUCCESS, nil)
		httpContext.JSON(http.StatusOK, response)
	}
}
