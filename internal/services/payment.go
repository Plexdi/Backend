package services

import (
	"fmt"
	"log"
	"strconv"

	"github.com/Plexdi/plexdi-studio-backend/internal/data"
	"github.com/stripe/stripe-go/v83"
	"github.com/stripe/stripe-go/v83/checkout/session"
)

func CreateCheckoutSession(priceId string, quantity int64, CommissionsID int64) (string, error) {
	params := &stripe.CheckoutSessionParams{
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			&stripe.CheckoutSessionLineItemParams{
				Price:    stripe.String(priceId),
				Quantity: stripe.Int64(quantity),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String("https://Plexdistudio.com/payments/success"),
		CancelURL:  stripe.String("https://Plexdistudio.com/payments/cancel"),

		AllowPromotionCodes: stripe.Bool(true),

		Metadata: map[string]string{
			"commissions_id": strconv.FormatInt(CommissionsID, 10),
		},
	}

	s, err := session.New(params)

	if err != nil {
		log.Printf("❌ Stripe Checkout error: %v\n", err)
		return "", err
	}

	return s.URL, nil
}

func FindProductTierPrice(product string, tier string) (string, error) {
	var tierIndex int
	priceList, ok := data.PriceMap[product]
	if !ok {
		return "", fmt.Errorf("unknown product: %s", product)
	}

	if product == "Discord Server Package" || product == "Discord User Package" ||
		product == "Social Media Banner Package" ||
		product == "Starter Streamer Pack Package" ||
		product == "Starter Youtube Package" ||
		product == "Streamer Package" {
		tierIndex = 0
		return priceList[tierIndex], nil
	}

	switch tier {
	case "Starter":
		tierIndex = 0
	case "Standard":
		tierIndex = 1
	case "Premium":
		tierIndex = 2
	default:
		return "", fmt.Errorf("unknown tier: %s", tier)
	}
	return priceList[tierIndex], nil
}
