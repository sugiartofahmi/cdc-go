package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"go-service/infrastructure/config"
	"go-service/infrastructure/databases"
	"go-service/infrastructure/integrations"
	"go-service/infrastructure/middlewares"
	"go-service/infrastructure/singleton"

	userInterfaces "go-service/domain/user/interfaces"
	userRepositories "go-service/domain/user/repositories"
	userServices "go-service/domain/user/services"

	roleInterfaces "go-service/domain/role/interfaces"
	roleRepositories "go-service/domain/role/repositories"
	roleServices "go-service/domain/role/services"

	authInterfaces "go-service/domain/auth/interfaces"
	authRepositories "go-service/domain/auth/repositories"
	authServices "go-service/domain/auth/services"

	healthController "go-service/presentation/api/health/controller"
	roleController "go-service/presentation/api/v1/role/controller"
	userController "go-service/presentation/api/v1/user/controller"
	authController "go-service/presentation/api/v1/auth/controller"

	redisFactory "go-service/infrastructure/redis/factories"
	redisInterfaces "go-service/infrastructure/redis/interfaces"
	redisServices "go-service/infrastructure/redis/services"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"

	"go-service/migration"
	"go-service/seeder"
)

var (
	router                       *gin.Engine
	redisCache                  redisInterfaces.RedisCacheInterface
	userQueryRepository         userInterfaces.UserQueryRepositoryInterface
	userStoreRepository         userInterfaces.UserStoreRepositoryInterface
	userRoleQueryRepository     userInterfaces.UserRoleQueryRepositoryInterface
	userService                 userInterfaces.UserServiceInterface
	roleQueryRepository         roleInterfaces.RoleQueryRepositoryInterface
	roleStoreRepository         roleInterfaces.RoleStoreRepositoryInterface
	roleService                 roleInterfaces.RoleServiceInterface
	authUserQueryRepository     authInterfaces.AuthUserQueryRepositoryInterface
	authUserStoreRepository     authInterfaces.AuthUserStoreRepositoryInterface
	authRoleQueryRepository     authInterfaces.AuthRoleQueryRepositoryInterface
	authService                 authInterfaces.AuthServiceInterface
	execMigration               *string
	flagMigration               *string
	migrationFileName           *string
	autoMigrateFlag             *string
	flagSeeder                 *string
	seederTarget               *string
)

func main() {
	extractArgs()
	initializeSingleton()
	runnerMigration()
	runnerAutoMigrate()
	runnerSeeder()
	initializeRouter()
	initializeRepositories()
	initializeServices()
	initializeControllers()
	initializeHttpServer()
}

func extractArgs() {
	execMigration = flag.String("exec", "up", "--exec [up/down/fresh/create]")
	flagMigration = flag.String("migration", "false", "--migration [true/false]")
	migrationFileName = flag.String("fileName", "", "--fileName <name>")
	autoMigrateFlag = flag.String("automigrate", "false", "--automigrate [true/false]")
	flagSeeder = flag.String("seed", "false", "--seed [true/false]")
	seederTarget = flag.String("target", "", "--target [SeederName,...] (optional)")
	flag.Parse()
}

func initializeSingleton() {
	db, err := databases.NewDBConnection()
	if err != nil {
		panic(err)
	}

	redis, err := redisFactory.NewRedisClient()
	if err != nil {
		panic(err)
	}

	var opensearch *opensearchapi.Client
	if config.OpensearchHost != "" && config.OpensearchPassword != "" {
		opensearch, err = databases.NewOpenSearchConnection()
		if err != nil {
			log.Printf("warn: opensearch connection failed: %v", err)
		}
	} else {
		log.Println("opensearch skipped: OPENSEARCH_HOST or OPENSEARCH_PASSWORD not set")
	}

	httpClient := integrations.NewHttpClient()

	singleton.Init(httpClient, db, redis, opensearch)
	redisCache = redisServices.NewRedisCacheService(redis)

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

func runnerMigration() {
	if *flagMigration != "true" {
		return
	}
	if *execMigration == "create" {
		migration.Create(nil, *migrationFileName)
		os.Exit(0)
	}
	migration.Run(singleton.PostgresSingleton(), *execMigration)
	os.Exit(0)
}

func runnerAutoMigrate() {
	if *autoMigrateFlag != "true" {
		return
	}
	if err := migration.AutoMigrate(singleton.PostgresSingleton()); err != nil {
		log.Fatal(err)
	}
	log.Println("automigrate completed")
	os.Exit(0)
}

func runnerSeeder() {
	if *flagSeeder != "true" {
		return
	}
	var classes []string
	if *seederTarget != "" {
		classes = strings.Split(*seederTarget, ",")
	}
	if err := seeder.Run(singleton.PostgresSingleton(), classes); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func initializeJobs() {
	log.Println("jobs initialized")
}

func initializeSchedulers() {
	log.Println("schedulers initialized")
}

func initializeRepositories() {
	userQueryRepository         = userRepositories.NewUserQueryRepository(singleton.PostgresSingleton())
	userStoreRepository         = userRepositories.NewUserStoreRepository(singleton.PostgresSingleton())
	userRoleQueryRepository     = userRepositories.NewUserRoleQueryRepository(singleton.PostgresSingleton())
	roleQueryRepository         = roleRepositories.NewRoleQueryRepository(singleton.PostgresSingleton())
	roleStoreRepository         = roleRepositories.NewRoleStoreRepository(singleton.PostgresSingleton())
	authUserQueryRepository     = authRepositories.NewAuthUserQueryRepository(singleton.PostgresSingleton())
	authUserStoreRepository     = authRepositories.NewAuthUserStoreRepository(singleton.PostgresSingleton())
	authRoleQueryRepository     = authRepositories.NewAuthRoleQueryRepository(singleton.PostgresSingleton())
	log.Println("repositories initialized")
}

func initializeServices() {
	userService = userServices.NewUserService(userQueryRepository, userStoreRepository, userRoleQueryRepository)
	roleService = roleServices.NewRoleService(roleQueryRepository, roleStoreRepository)
	authService = authServices.NewAuthService(authUserQueryRepository, authUserStoreRepository, authRoleQueryRepository)
	log.Println("services initialized")
}

func initializeControllers() {
	healthController.NewHealthController(router)
	roleController.NewRoleController(router, roleService)
	userController.NewUserController(router, userService)
	authController.NewAuthController(router, authService)
	log.Println("controllers initialized")
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
