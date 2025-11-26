package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Plexdi/plexdi-studio-backend/internal/db"
	"github.com/Plexdi/plexdi-studio-backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v83"
	"github.com/stripe/stripe-go/v83/webhook"
)

// ---------------------- Types ---------------------------

type PaymentRequest struct {
	Product      string `json:"item"`
	Amount       int64  `json:"amount"`
	CommissionID int64  `json:"commissionId"`
	Tier         string `json:"tier"`
	Discount     string `json:"discountCode"`
}

// ---------------------- routes registerations ---------------------------

func RegisterPaymentRoutes(r *gin.Engine) {
	r.POST("/payments/createCheckoutSession", createCheckoutSession)
	r.POST("/payments/webhook", stripeWebhook)

}

// ---------------------- controllers ---------------------------

func createCheckoutSession(c *gin.Context) {
	var req PaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("❌ createCheckoutSession: invalid request:", err)
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	priceID, err := services.FindProductTierPrice(req.Product, req.Tier)
	if err != nil {
		log.Println("❌ createCheckoutSession: price lookup failed:", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Create the checkout session
	url, err := services.CreateCheckoutSession(priceID, req.Amount, req.CommissionID)
	if err != nil {
		log.Println("❌ createCheckoutSession: failed to create session:", err)
		c.JSON(500, gin.H{"error": "failed to create checkout session", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"url": url})

}

func stripeWebhook(c *gin.Context) {
	var comm services.Commission
	//1) read the body
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Status(http.StatusBadRequest)
	}

	//get stripe signature header + secret
	sigHeader := c.GetHeader("Stripe-Signature")
	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")

	// 3) Verify request and construct event
	event, err := webhook.ConstructEvent(payload, sigHeader, endpointSecret)
	if err != nil {
		log.Println("❌ webhook: signature verification failed:", err)
		c.Status(http.StatusBadRequest)
		return
	}

	// 4) Handle the event
	switch event.Type {
	case "checkout.session.completed":
		var sess stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &sess); err != nil {
			log.Println("❌ webhook: failed to parse session:", err)
			c.Status(http.StatusBadRequest)
			return
		}

		// 5) Get your commission ID from metadata
		commissionIDStr := sess.Metadata["commission_id"]
		if commissionIDStr == "" {
			log.Println("❌ webhook: missing commission_id in metadata")
			c.Status(http.StatusBadRequest)
			return
		}

		commissionID, err := strconv.Atoi(commissionIDStr)
		if err != nil {
			log.Println("❌ webhook: invalid commission_id:", commissionIDStr, err)
			c.Status(http.StatusBadRequest)
			return
		}

		log.Println("✅ webhook: payment completed for commission", commissionID)

		// 6) Update DB: mark commission as paid
		_, err = db.Pool.Exec(
			context.Background(),
			`UPDATE commissions SET status = $1 WHERE id = $2`,
			"paid",
			commissionID,
		)
		if err != nil {
			log.Println("❌ webhook: failed to update commission status:", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		err = db.Pool.QueryRow(
			context.Background(),
			`SELECT id, name, email, discord, type, details, status
			FROM commissions
			WHERE id = $1`,
			commissionID,
		).Scan(
			&comm.ID,
			&comm.Name,
			&comm.Email,
			&comm.Discord,
			&comm.Type,
			&comm.Details,
			&comm.Status,
		)

		if err != nil {
			log.Println("❌ webhook: failed to load commission for email:", err)
			c.Status(http.StatusInternalServerError)
			return
		}

	default:
		// ignore other events or log them
		log.Println("ℹ️ webhook: ignoring event type", event.Type)

	}

	c.Status(http.StatusOK)
}
