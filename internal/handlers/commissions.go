package handlers

import (
	"net/http"

	"github.com/Plexdi/plexdi-studio-backend/internal/services"
	"github.com/gin-gonic/gin"
)

func RegisterCommissionRoutes(r *gin.Engine) {
	r.POST("/commissions", CreateCommission)
	r.GET("/commissions", GetAllCommissions)
}

func CreateCommission(c *gin.Context) {
	var req services.Commission
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	newCommission := services.MakeCommission(req.Name, req.Email, req.Type, req.Details)
	c.JSON(http.StatusCreated, newCommission)

	services.SendCommissionEmail(req.Email, services.CommissionData{
		Name: req.Name,
		Type: req.Type,
	})

}

func GetAllCommissions(c *gin.Context) {
	all := services.GetAllCommissions()
	c.JSON(http.StatusOK, all)
}
