package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"go-service/infrastructure/config"
	"go-service/infrastructure/databases"
	"go-service/infrastructure/integrations"
	"go-service/infrastructure/middlewares"
	"go-service/infrastructure/singleton"

	kafkaFactory "go-service/infrastructure/kafka/factories"
	redisFactory "go-service/infrastructure/redis/factories"

	productInterfaces "go-service/domain/product/interfaces"
	productRepositories "go-service/domain/product/repositories"
	productServices "go-service/domain/product/services"

	healthController "go-service/presentation/api/health/controller"
	productController "go-service/presentation/api/v1/product/controller"
)

var (
	router                 *gin.Engine
	productQueryRepository productInterfaces.ProductQueryRepositoryInterface
	productService         productInterfaces.ProductServiceInterface
)

func main() {
	initializeSingleton()
	initializeRouter()
	initializeRepositories()
	initializeServices()
	initializeControllers()
	initializeIndexer()
	initializeHttpServer()
}

func initializeSingleton() {
	redis, err := redisFactory.NewRedisClient()
	if err != nil {
		panic(err)
	}

	opensearch, err := databases.NewOpenSearchConnection()
	if err != nil {
		panic(err)
	}

	httpClient := integrations.NewHttpClient()

	kafka, err := kafkaFactory.NewKafka()
	if err != nil {
		panic(err)
	}

	singleton.Init(httpClient, redis, opensearch, kafka)

	log.Println("singletons initialized")
}

func initializeRouter() {
	router = gin.New()
	router.ContextWithFallback = true

	gin.SetMode(config.AppGinMode)

	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	router.Use(gin.Recovery())
	router.Use(middlewares.LoggingMiddleware())
	router.Use(cors.New(corsConfig))
	router.Use(middlewares.ExceptionMiddleware())
}

func initializeRepositories() {
	productQueryRepository = productRepositories.NewProductQueryRepository(singleton.OpensearchSingleton())
	log.Println("repositories initialized")
}

func initializeServices() {
	productService = productServices.NewProductService(productQueryRepository)
	log.Println("services initialized")
}

func initializeControllers() {
	healthController.NewHealthController(router)
	productController.NewProductController(router, productService)
	log.Println("controllers initialized")
}

func initializeIndexer() {
	log.Println("indexer: placeholder — will be implemented in Phase 3")
}

func initializeHttpServer() {
	srv := &http.Server{
		Addr:    ":" + config.AppPort,
		Handler: router,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("starting server on port %s", config.AppPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("server stopped")
}
