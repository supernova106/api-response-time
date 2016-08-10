package main

import (
	"net/http"
	"os"

	"api-response-time/app/db"
	"api-response-time/app/handlers"
	"github.com/gin-gonic/gin"
	"log"
)

const (
	// Port at which the server starts listening
	Port = "8080"
)

var mdb *db.DB

func main() {
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	var err error
	mdb, err = db.Connect()
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	router := gin.Default()

	router.Use(injectDependencyServices())
	router.GET("/api", accesslog.List)

	// This handler will match /user/john but will not match neither /user/ or /user
	router.GET("/api/:apiname", func(c *gin.Context) {
		apiname := c.Param("apiname")
		c.String(http.StatusOK, "Hello %s", apiname)
	})

	// By default it serves on :8080 unless a
	// API_PORT environm+nt variable was defined.
	port := os.Getenv("API_PORT")

	if len(port) == 0 {
		port = Port
	}

	router.Run(":" + port)
	// router.Run(":3000") for a hard coded port
}

func injectDependencyServices() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("mdb", mdb)
		c.Next()
	}
}
