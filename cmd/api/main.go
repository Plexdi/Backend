package main

import (
	"log"
	"os"

	"github.com/Plexdi/plexdi-studio-backend/internal/data"
	"github.com/Plexdi/plexdi-studio-backend/internal/db"
	"github.com/Plexdi/plexdi-studio-backend/internal/handlers"
	"github.com/Plexdi/plexdi-studio-backend/internal/middleware"
	"github.com/Plexdi/plexdi-studio-backend/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v83"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found — using Render environment variables")
	}

	if err := db.ConnectDB(); err != nil {
		log.Fatal(err)
	}

	defer db.Pool.Close()

	data.InitPriceMap()

	sk := os.Getenv("Stripe_Secret_key")
	if sk == "" {
		log.Fatal("❌ Stripe_Secret_key not set")
	} else {
		log.Println("✅ Stripe API key loaded")
	}

	stripe.Key = sk

	services.LoadCommissions()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://plexdistudio.com", "http://localhost:10000", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	r.Use(middleware.LimitRequests())

	handlers.RegisterCommissionRoutes(r)
	handlers.RegisterPaymentRoutes(r)

	r.Run(":" + os.Getenv("PORT"))
}
