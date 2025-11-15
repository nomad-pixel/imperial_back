package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nomad-pixel/imperial/pkg/postgres"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pgUrl := os.Getenv("DATABASE_URL")
	if pgUrl == "" {
		log.Fatalf("DATABASE_URL is not set")
	}

	db, err := postgres.NewPool(ctx, pgUrl)
	if err != nil {
		log.Fatalf("failed to create database pool: %v", err)
	}
	defer db.Close()

	server := gin.Default()

	server.Run(":8080")

}
