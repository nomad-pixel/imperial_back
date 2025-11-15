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

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pgUrl := os.Getenv("DATABASE_URL")
	if pgUrl == "" {
		log.Fatalf("DATABASE_URL is not set")
	}

	// Инициализация всех зависимостей через DI контейнер
	app, err := di.InitializeApp(ctx, pgUrl)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}
	defer app.Close()

	// Настройка HTTP сервера
	server := gin.New()

	// Middleware
	server.Use(gin.Logger())
	server.Use(middleware.Recovery())
	server.Use(middleware.ErrorHandler())

	// Swagger документация
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Создание группы /api для всех маршрутов
	apiGroup := server.Group("/api")

	// Регистрация маршрутов
	auth.RegisterRoutes(apiGroup, app.AuthHandler)
	if err := server.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
