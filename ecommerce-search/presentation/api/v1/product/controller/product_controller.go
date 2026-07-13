package controller

import (
	"net/http"

	productConstants "go-service/domain/product/constants"
	productDtos "go-service/domain/product/dtos"
	productInterfaces "go-service/domain/product/interfaces"
	"go-service/infrastructure/utils"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	productService productInterfaces.ProductServiceInterface
}

func NewProductController(router *gin.Engine, productService productInterfaces.ProductServiceInterface) {
	controller := &ProductController{productService: productService}
	router.GET("/products", controller.Pagination())
}

func (c *ProductController) Pagination() gin.HandlerFunc {
	return func(httpContext *gin.Context) {
		ctx := httpContext.Request.Context()
		dto := productDtos.AssignProductQueryRequestDto(httpContext)
		utils.Paginate(&dto.PaginationQueryRequestDto)
		result := c.productService.Pagination(ctx, dto)
		meta := utils.PaginationMetaBuilder(dto.Page, dto.PerPage, int(result.Count))
		response := utils.SuccessResponsePagination(http.StatusOK, productConstants.PRODUCT_PAGINATION_SUCCESS, result.Data, *meta)
		httpContext.JSON(http.StatusOK, response)
	}
}
