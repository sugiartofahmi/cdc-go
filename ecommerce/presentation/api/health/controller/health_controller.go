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
	Service  string `json:"service"`
	Database string `json:"database"`
	Redis    string `json:"redis"`
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
			Service:  "ok",
			Database: "ok",
			Redis:    "ok",
		}

		sqlDB, err := singleton.PostgresSingleton().DB()
		if err != nil {
			status.Database = "error: " + err.Error()
		} else if err := sqlDB.PingContext(ctx); err != nil {
			status.Database = "error: " + err.Error()
		}

		redisClient := singleton.RedisSingleton()
		if redisClient == nil {
			status.Redis = "error: redis client not initialized"
		} else if _, err := redisClient.Ping(ctx).Result(); err != nil {
			status.Redis = "error: " + err.Error()
		}

		allHealthy := status.Service == "ok" && status.Database == "ok" && status.Redis == "ok"

		if allHealthy {
			response := utils.SuccessResponse(http.StatusOK, "service is healthy", status)
			httpContext.JSON(http.StatusOK, response)
		} else {
			response := utils.SuccessResponse(http.StatusServiceUnavailable, "service is unhealthy", status)
			httpContext.JSON(http.StatusServiceUnavailable, response)
		}
	}
}
