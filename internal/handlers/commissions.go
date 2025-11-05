package handlers

import (
	"context"
	"net/http"

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
		`INSERT INTO commissions (id, name, email, discord, type, details)
		 VALUES (gen_random_uuid(), $1, $2, $3, $4, $5)`,
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
	go services.SendCommissionEmail(req.Email, services.CommissionData{
		Name: req.Name,
		Type: req.Type,
	})

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
