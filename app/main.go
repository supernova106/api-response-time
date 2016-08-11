package main

import (
	"api-response-time/app/config"
	"api-response-time/app/db"
	"api-response-time/app/handlers"
	"api-response-time/app/middlewares"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// Load config
	cfg, err := config.Load(".env")
	if err != nil {
		log.Fatalf("Can't load .env file %v", err)
		return
	}

	// connect DB
	db.Connect(cfg.MongoDBUrl)

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware

	router := gin.Default()

	// Middlewares
	router.Use(middlewares.Connect)
	router.Use(middlewares.ErrorHandler)

	router.GET("/api", accesslog.GetAll)

	// By default it serves on :8080 unless a
	// API_PORT environm+nt variable was defined.
	router.Run(":" + cfg.Port)
}
