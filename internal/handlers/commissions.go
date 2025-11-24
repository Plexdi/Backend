package handlers

import (
	"context"
	"net/http"
	"strconv"

	"log"

	"github.com/Plexdi/plexdi-studio-backend/internal/db"
	"github.com/Plexdi/plexdi-studio-backend/internal/services"
	"github.com/gin-gonic/gin"
)

// ---------------------- routes registerations ---------------------------

func RegisterCommissionRoutes(r *gin.Engine) {
	r.POST("/commissions", CreateCommission)
	r.GET("/commissions", GetAllCommissions)
	r.PATCH("/commissions/:id", UpdateCommissions)
	r.DELETE("/commissions/:id", DeleteCommission)
}

// ---------------------- controllers ---------------------------

func CreateCommission(c *gin.Context) {
	var req services.Commission
	var newID int64

	// Parse JSON body
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Save to (PostgreSQL)
	err := db.Pool.QueryRow(
		context.Background(),
		`INSERT INTO commissions (name, email, discord, type, details, status, designers)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`,
		req.Name, req.Email, req.Discord, req.Type, req.Details, "queued", req.Designers,
	).Scan(&newID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "database insert failed",
			"details": err.Error(),
		})
		return
	}

	// Send response
	c.JSON(http.StatusCreated, gin.H{
		"message": "Form submitted successfully. Please check your email.",
		"id":      newID,
	})

}

func GetAllCommissions(c *gin.Context) {
	data, err := services.GetAllCommissions()
	if err != nil {
		log.Printf("❌ error fetching commissions: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch commissions",
		})
		return
	}

	c.JSON(http.StatusOK, data)
}

func UpdateCommissions(c *gin.Context) {
	// 1. ID from URL
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid commission ID"})
		return
	}

	// 2. JSON body
	var req struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON body"})
		return
	}

	if req.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing status"})
		return
	}

	// 3. Call service
	if err := services.UpdateCommissionStatus(id, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update commission",
			"error":   err.Error(), // 👈 TEMP: show real error
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Commission updated successfully",
		"id":      id,
		"status":  req.Status,
	})
}

func DeleteCommission(c *gin.Context) {
	// 1. ID from URL
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid commission ID"})
		return
	}
	// 2. Call service
	if err := services.DeleteCommission(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete commission",
			"error":   err.Error(), // 👈 TEMP: show real error
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Commission deleted successfully",
		"id":      id,
	})

}
