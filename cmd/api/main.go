package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/nomad-pixel/imperial/docs"
	"github.com/nomad-pixel/imperial/internal/di"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/auth"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/car"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/car_category"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/car_image"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/car_mark"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/car_tag"
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

	pgUrl := os.Getenv("DATABASE_URL")
	if pgUrl == "" {
		log.Fatalf("DATABASE_URL is not set")
	}

	app, err := di.InitializeApp(ctx, pgUrl)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}
	defer app.Close()

	server := gin.New()

	server.Use(gin.Logger())
	server.Use(middleware.Recovery())
	server.Use(middleware.ErrorHandler())

	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server.Static("/images", "./uploads")

	apiGroup := server.Group("/api")

	auth.RegisterRoutes(apiGroup, app.AuthHandler)
	car.RegisterRoutes(apiGroup, app.CarHandler, app.TokenService)
	car_tag.RegisterRoutes(apiGroup, app.CarTagHandler, app.TokenService)
	car_mark.RegisterRoutes(apiGroup, app.CarMarkHandler, app.TokenService)
	car_category.RegisterRoutes(apiGroup, app.CarCategoryHandler, app.TokenService)
	car_image.RegisterRoutes(apiGroup, app.CarImageHandler, app.TokenService)
	if err := server.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
