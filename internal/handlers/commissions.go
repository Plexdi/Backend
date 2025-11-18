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

	// Parse JSON body
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Save to (PostgreSQL)
	_, err := db.Conn.Exec(context.Background(),
		`INSERT INTO commissions (name, email, discord, type, details, status, designers)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,

		req.Name, req.Email, req.Discord, req.Type, req.Details, "queued", req.Designers,
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
	go func() {
		log.Println("🚀 SendCommissionEmail() triggered")

		err := services.SendCommissionEmail(req.Email, services.Commission{
			Name:    req.Name,
			Email:   req.Email,
			Discord: req.Discord,
			Type:    req.Type,
		})
		if err != nil {
			log.Printf("❌ Email failed for %s: %v\n", req.Email, err)
		} else {
			log.Printf("✅ Email sent successfully to %s\n", req.Email)
		}
	}()

	// Send response
	c.JSON(http.StatusCreated, gin.H{
		"message":    "Form submitted successfully. Please check your email.",
		"commission": newCommission,
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
