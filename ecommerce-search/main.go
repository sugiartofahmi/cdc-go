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
	"github.com/segmentio/kafka-go"

	"go-service/infrastructure/config"
	"go-service/infrastructure/databases"
	"go-service/infrastructure/integrations"
	"go-service/infrastructure/middlewares"
	"go-service/infrastructure/singleton"

	kafkaFactory "go-service/infrastructure/kafka/factories"
	kafkaInterfaces "go-service/infrastructure/kafka/interfaces"
	kafkaServices "go-service/infrastructure/kafka/services"
	redisFactory "go-service/infrastructure/redis/factories"

	productDtos "go-service/domain/product/dtos"
	productInterfaces "go-service/domain/product/interfaces"
	productRepositories "go-service/domain/product/repositories"
	productServices "go-service/domain/product/services"

	healthController "go-service/presentation/api/health/controller"
	productController "go-service/presentation/api/v1/product/controller"
)

var (
	router                 *gin.Engine
	productQueryRepository productInterfaces.ProductQueryRepositoryInterface
	productStoreRepository productInterfaces.ProductStoreRepositoryInterface
	productService         productInterfaces.ProductServiceInterface
	productEventService    productInterfaces.ProductEventServiceInterface
	kafkaConsumer          kafkaInterfaces.KafkaConsumerInterface
	consumerShutdown       context.CancelFunc
)

func main() {
	initializeSingleton()
	initializeRouter()
	initializeRepositories()
	initializeServices()
	initializeControllers()
	initializeKafkaConsumer()
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
	productStoreRepository = productRepositories.NewProductStoreRepository(singleton.OpensearchSingleton())
	log.Println("repositories initialized")
}

func initializeServices() {
	productService = productServices.NewProductService(productQueryRepository)
	productEventService = productServices.NewProductEventService(productStoreRepository)
	log.Println("services initialized")
}

func initializeControllers() {
	healthController.NewHealthController(router)
	productController.NewProductController(router, productService)
	log.Println("controllers initialized")
}

func initializeKafkaConsumer() {
	ctx, cancel := context.WithCancel(context.Background())
	consumerShutdown = cancel

	kafkaConsumer = kafkaServices.NewKafkaConsumerService(singleton.KafkaSingleton())

	handler := func(msg kafka.Message) error {
		switch msg.Topic {
		case "ecommerce.public.products":
			event := productDtos.ProductEventDtoFromMessage(msg)
			if event != nil {
				productEventService.Handle(ctx, event)
			}
		}
		return nil
	}

	go func() {
		log.Println("kafka consumer started")
		kafkaConsumer.Consume(ctx, handler)
		log.Println("kafka consumer stopped")
	}()
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

	if consumerShutdown != nil {
		consumerShutdown()
	}
	if kafkaConsumer != nil {
		if err := kafkaConsumer.Close(); err != nil {
			log.Println("error closing kafka consumer:", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("server stopped")
}
