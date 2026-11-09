package api

import (
	"fmt"
	"net/http"

	_ "github.com/Baroneee/EBVN_Assignment_1/docs"
	"github.com/Baroneee/EBVN_Assignment_1/internal/handler"
	"github.com/Baroneee/EBVN_Assignment_1/internal/repository"
	"github.com/Baroneee/EBVN_Assignment_1/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Engine starts the HTTP server and handles HTTP requests.
type Engine interface {
	Start() error
	ServerHTTP(w http.ResponseWriter, r *http.Request)
}

// engine is the Gin-based implementation of Engine.
type engine struct {
	app         *gin.Engine
	cfg         *Config
	redisClient *redis.Client
}

// NewEngine creates an HTTP engine configured with cfg.
func NewEngine(cfg *Config, redisClient *redis.Client) Engine {
	e := &engine{
		app:         gin.Default(),
		cfg:         cfg,
		redisClient: redisClient,
	}
	e.initRoutes()
	return e
}

// Start runs the HTTP server on the configured application port.
func (e *engine) Start() error {
	return e.app.Run(fmt.Sprintf(":%d", e.cfg.AppPort))
}

// ServerHTTP serves an HTTP request through the configured Gin engine.
func (e *engine) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	e.app.ServeHTTP(w, r)
}

// initRoutes registers the application's HTTP routes.
func (e *engine) initRoutes() {
	healthCheckSvc := service.NewHealthCheck(
		e.cfg.ServiceName,
		e.cfg.InstanceID,
	)

	urlStorageRepo := repository.NewURLStorage(e.redisClient)
	shortenURLSvc := service.NewShortenURL(urlStorageRepo)
	shortenURLHandler := handler.NewShortenURL(shortenURLSvc)

	healthCheckHandler := handler.NewHealthCheck(healthCheckSvc)
	e.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	e.app.GET("/health", healthCheckHandler.Check)
	e.app.POST("/v1/links/shortern", shortenURLHandler.CreateShorternLink)
}
