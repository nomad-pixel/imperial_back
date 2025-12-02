package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/nomad-pixel/imperial/docs"
	"github.com/nomad-pixel/imperial/internal/di"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/auth"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/car"
	carCategory "github.com/nomad-pixel/imperial/internal/interfaces/http/car/category"
	carImage "github.com/nomad-pixel/imperial/internal/interfaces/http/car/image"
	carMark "github.com/nomad-pixel/imperial/internal/interfaces/http/car/mark"
	carTag "github.com/nomad-pixel/imperial/internal/interfaces/http/car/tag"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/celebrity"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/middleware"
)

// @title           Imperial API
// @version         1.0
// @description     REST API для проекта Imperial
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@imperial.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
// @description JWT token must be passed with `Bearer ` prefix. Example: "Bearer eyJhbGciOiJIUzI1NiI..."

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize application with all dependencies
	app, err := di.InitializeApp(ctx)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}
	defer app.Close()

	cfg := app.Config

	log.Printf("🚀 Starting %s server in %s mode", cfg.App.Name, cfg.App.Environment)
	if cfg.IsDevelopment() {
		log.Println("⚠️  Debug mode enabled")
	}

	server := gin.New()

	server.Use(gin.Logger())
	server.Use(middleware.Recovery())
	server.Use(middleware.ErrorHandler())

	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Serve static files from storage path
	server.Static("/uploads", cfg.Storage.LocalPath)

	apiGroup := server.Group("/api")

	auth.RegisterRoutes(apiGroup, app.AuthHandler)
	car.RegisterRoutes(apiGroup, app.CarHandler, app.TokenService)
	carTag.RegisterRoutes(apiGroup, app.CarTagHandler, app.TokenService)
	carMark.RegisterRoutes(apiGroup, app.CarMarkHandler, app.TokenService)
	carCategory.RegisterRoutes(apiGroup, app.CarCategoryHandler, app.TokenService)
	carImage.RegisterRoutes(apiGroup, app.CarImageHandler, app.TokenService)
	celebrity.RegisterRoutes(apiGroup, app.CelebrityHandler, app.TokenService)

	log.Printf("✅ Server listening on http://localhost:%d", cfg.Server.Port)
	log.Printf("📚 Swagger documentation: http://localhost:%d/swagger/index.html", cfg.Server.Port)

	if err := server.Run(cfg.GetServerAddress()); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
