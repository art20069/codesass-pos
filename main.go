package main

import (
	"codezard-pos/db"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	db.ConnectDB()
	db.Migrate()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true

	os.MkdirAll("uploads/products", 0755)
	r := gin.Default()
	r.Use(cors.New(corsConfig))
	r.Static("/uploads", "./uploads")
	serveRoutes(r)

	port := os.Getenv("PORT")
	if port != "" {
		port = "5000"
	}
	r.Run(":" + port)
}
