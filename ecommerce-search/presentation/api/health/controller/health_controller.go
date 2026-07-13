package controller

import (
	"context"
	"net/http"
	"time"

	"go-service/infrastructure/singleton"
	"go-service/infrastructure/utils"

	"github.com/gin-gonic/gin"
)

type HealthStatus struct {
	Service    string `json:"service"`
	Redis      string `json:"redis"`
	OpenSearch string `json:"opensearch"`
}

type HealthController struct{}

func NewHealthController(router *gin.Engine) {
	controller := &HealthController{}
	router.GET("/health", controller.Check())
}

func (h *HealthController) Check() gin.HandlerFunc {
	return func(httpContext *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		status := HealthStatus{
			Service:    "ok",
			Redis:      "ok",
			OpenSearch: "ok",
		}

		redisClient := singleton.RedisSingleton()
		if redisClient == nil {
			status.Redis = "error: redis client not initialized"
		} else if _, err := redisClient.Ping(ctx).Result(); err != nil {
			status.Redis = "error: " + err.Error()
		}

		opensearchClient := singleton.OpensearchSingleton()
		if opensearchClient == nil {
			status.OpenSearch = "error: opensearch client not initialized"
		} else if _, err := opensearchClient.Ping(ctx, nil); err != nil {
			status.OpenSearch = "error: " + err.Error()
		}

		allHealthy := status.Service == "ok" && status.Redis == "ok" && status.OpenSearch == "ok"

		if allHealthy {
			response := utils.SuccessResponse(http.StatusOK, "service is healthy", status)
			httpContext.JSON(http.StatusOK, response)
		} else {
			response := utils.SuccessResponse(http.StatusServiceUnavailable, "service is unhealthy", status)
			httpContext.JSON(http.StatusServiceUnavailable, response)
		}
	}
}
