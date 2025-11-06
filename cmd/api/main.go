package main

import (
	"log"
	"os"

	"github.com/Plexdi/plexdi-studio-backend/internal/db"
	"github.com/Plexdi/plexdi-studio-backend/internal/handlers"
	"github.com/Plexdi/plexdi-studio-backend/internal/middleware"
	"github.com/Plexdi/plexdi-studio-backend/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found — using Render environment variables")
	}

	if err := db.ConnectDB(); err != nil {
		log.Fatal(err)
	}

	services.LoadCommissions()

	r := gin.Default()
	r.Use(middleware.LimitRequests())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))
	handlers.RegisterCommissionRoutes(r)
	r.Run(":" + os.Getenv("PORT"))
}
