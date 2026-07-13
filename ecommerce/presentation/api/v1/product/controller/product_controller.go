package controller

import (
	"net/http"

	productConstants "go-service/domain/product/constants"
	productDtos "go-service/domain/product/dtos"
	productInterfaces "go-service/domain/product/interfaces"
	"go-service/infrastructure/middlewares"
	"go-service/infrastructure/utils"
	"go-service/infrastructure/validators"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	productService productInterfaces.ProductServiceInterface
}

func NewProductController(router *gin.Engine, productService productInterfaces.ProductServiceInterface) {
	controller := &ProductController{productService: productService}

	productRoute := router.Group("/api/v1/products")
	productRoute.GET("", controller.Pagination())
	productRoute.GET("/:id", controller.Detail())
	productRoute.POST("", middlewares.ValidateRequestJson[productDtos.ProductCreateRequestDto](), controller.Create())
	productRoute.PUT("/:id", middlewares.ValidateRequestJson[productDtos.ProductUpdateRequestDto](), controller.Update())
	productRoute.DELETE("/:id", controller.Delete())
}

func (c *ProductController) Pagination() gin.HandlerFunc {
	return func(httpContext *gin.Context) {
		ctx := httpContext.Request.Context()
		dto := productDtos.AssignProductQueryRequestDto(httpContext)
		result := c.productService.Pagination(ctx, dto)
		meta := utils.PaginationMetaBuilder(dto.Page, dto.PerPage, int(result.Count))
		items := productDtos.ProductResponseDtoFromEntities(result.Data)
		response := utils.SuccessResponsePagination(http.StatusOK, productConstants.PRODUCT_PAGINATION_SUCCESS, items, *meta)
		httpContext.JSON(http.StatusOK, response)
	}
}

func (c *ProductController) Detail() gin.HandlerFunc {
	return func(httpContext *gin.Context) {
		ctx := httpContext.Request.Context()
		id := validators.ValidateUUID(httpContext.Param("id"))
		result := productDtos.ProductResponseDtoFromEntity(c.productService.Detail(ctx, id))
		response := utils.SuccessResponse(http.StatusOK, productConstants.PRODUCT_DETAIL_SUCCESS, result)
		httpContext.JSON(http.StatusOK, response)
	}
}

func (c *ProductController) Create() gin.HandlerFunc {
	return func(httpContext *gin.Context) {
		ctx := httpContext.Request.Context()
		dto := httpContext.MustGet(middlewares.RequestBodyJsonKey).(*productDtos.ProductCreateRequestDto)
		result := productDtos.ProductResponseDtoFromEntity(c.productService.Create(ctx, dto))
		response := utils.SuccessResponse(http.StatusCreated, productConstants.PRODUCT_CREATE_SUCCESS, result)
		httpContext.JSON(http.StatusCreated, response)
	}
}

func (c *ProductController) Update() gin.HandlerFunc {
	return func(httpContext *gin.Context) {
		ctx := httpContext.Request.Context()
		id := validators.ValidateUUID(httpContext.Param("id"))
		dto := httpContext.MustGet(middlewares.RequestBodyJsonKey).(*productDtos.ProductUpdateRequestDto)
		result := productDtos.ProductResponseDtoFromEntity(c.productService.Update(ctx, id, dto))
		response := utils.SuccessResponse(http.StatusOK, productConstants.PRODUCT_UPDATE_SUCCESS, result)
		httpContext.JSON(http.StatusOK, response)
	}
}

func (c *ProductController) Delete() gin.HandlerFunc {
	return func(httpContext *gin.Context) {
		ctx := httpContext.Request.Context()
		id := validators.ValidateUUID(httpContext.Param("id"))
		c.productService.Delete(ctx, id)
		response := utils.SuccessResponse(http.StatusOK, productConstants.PRODUCT_DELETE_SUCCESS, nil)
		httpContext.JSON(http.StatusOK, response)
	}
}
