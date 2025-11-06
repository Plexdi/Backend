package handlers

import (
	"context"
	"net/http"

	"log"

	"github.com/Plexdi/plexdi-studio-backend/internal/db"
	"github.com/Plexdi/plexdi-studio-backend/internal/services"
	"github.com/gin-gonic/gin"
)

func RegisterCommissionRoutes(r *gin.Engine) {
	r.POST("/commissions", CreateCommission)
	r.GET("/commissions", GetAllCommissions)
}

func CreateCommission(c *gin.Context) {
	var req services.Commission

	// Parse JSON body
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Save to Supabase (PostgreSQL)
	_, err := db.Conn.Exec(context.Background(),
		`INSERT INTO commissions (name, email, discord, type, details)
		VALUES ($1, $2, $3, $4, $5)`,
		req.Name, req.Email, req.Discord, req.Type, req.Details,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "database insert failed",
			"details": err.Error(),
		})
		return
	}

	// Generate a local struct for the frontend
	newCommission := services.MakeCommission(req.Name, req.Email, req.Type, req.Details)

	// Send confirmation email asynchronously
	// Send confirmation email asynchronously
	go func() {
		log.Println("🚀 SendCommissionEmail() triggered")

		err := services.SendCommissionEmail(req.Email, services.CommissionData{
			Name: req.Name,
			Type: req.Type,
		})
		if err != nil {
			log.Printf("❌ Email failed for %s: %v\n", req.Email, err)
		} else {
			log.Printf("✅ Email sent successfully to %s\n", req.Email)
		}
	}()

	// Send response
	c.JSON(http.StatusCreated, gin.H{
		"message":    "✅ Commission created successfully!",
		"commission": newCommission,
	})
}

func GetAllCommissions(c *gin.Context) {
	all := services.GetAllCommissions()
	c.JSON(http.StatusOK, all)
}
