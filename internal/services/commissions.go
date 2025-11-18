package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"bytes"

	"github.com/Plexdi/plexdi-studio-backend/internal/db"
)

type Commission struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Discord   string    `json:"discord"`
	Details   string    `json:"details"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type CommissionData struct {
	Name string
	Type string
}

type resendEmail struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	HTML    string `json:"html"`
}

// ---------------------------- public variables --------------------------

var commissions []Commission
var nextID = 1

const filePath = "data/commissions.json"

// ----------------------------- make commmissions ------------------------------

func MakeCommission(name, email, ctype, details string) Commission {
	newCommission := Commission{
		ID:     nextID,
		Name:   name,
		Email:  email,
		Type:   ctype,
		Status: "queued",
	}
	nextID++
	commissions = append(commissions, newCommission)
	saveCommissions()
	fmt.Println("Commission created:", newCommission)
	return newCommission
}

// ------------------- save commissions into json file --------------------------
func saveCommissions() {
	if err := os.MkdirAll("data", os.ModePerm); err != nil {
		fmt.Println("Error creating data folder:", err)
		return
	}

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// IMPORTANT: write the slice, not nil or uninitialized variable
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // makes it pretty
	if err := encoder.Encode(commissions); err != nil {
		fmt.Println("Error writing commissions:", err)
	} else {
		fmt.Println("Saved", len(commissions), "commissions to file.")
	}
}

// ------------------- Load commissions --------------------------

func LoadCommissions() {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("No existing file, starting fresh.")
		return
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&commissions); err != nil {
		fmt.Println("Error decoding commissions:", err)
		return
	}

	if len(commissions) > 0 {
		nextID = commissions[len(commissions)-1].ID + 1
	}

	fmt.Println("Loaded", len(commissions), "commissions from file.")
}

func GetAllCommissions() ([]Commission, error) {

	rows, err := db.Conn.Query(context.Background(), `
		SELECT id, name, email, discord, details, type, status, created_at
		FROM commissions
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commissions []Commission

	for rows.Next() {
		var c Commission
		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Email,
			&c.Discord,
			&c.Details,
			&c.Type,
			&c.Status,
			&c.CreatedAt,
		); err != nil {
			return nil, err
		}
		commissions = append(commissions, c)
	}

	return commissions, nil
}

// ------------------ Send -------------------------------------------
func SendCommissionEmail(to string, data Commission) error {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		log.Println("⚠️ RESEND_API_KEY not found in environment")
		return nil
	}

	// =============== CLIENT EMAIL ===============
	clientHTML := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: auto; padding: 25px; background-color: #f8f9fa; border-radius: 10px;">
			<h2 style="color: #2563eb; text-align: center;">🎨 New Commission Received</h2>
			<p style="font-size: 16px; color: #333;">
				Hey <b>%s</b>,<br><br>
				We’ve received your commission request for a <b>%s</b>.<br>
				Our team will review it and get back to you shortly with details and payment options.
			</p>
			<hr style="border: none; border-top: 1px solid #ddd; margin: 25px 0;">
			<p style="text-align: center; font-size: 14px; color: #555;">
				– The Plexdi Studio Team<br>
				<a href="https://plexdistudio.com" style="color: #2563eb; text-decoration: none;">plexdistudio.com</a>
			</p>
		</div>`,
		data.Name, data.Type,
	)

	// =============== ADMIN EMAIL (NOTIFICATION) ===============
	adminHTML := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; background-color: #f4f7fa; padding: 25px; border-radius: 10px; max-width: 600px; margin: auto;">
			<h2 style="color: #2563eb; text-align: center;">🆕 New Commission Request</h2>
			<p style="font-size: 16px; color: #333;">
				A new commission has just been submitted through your Plexdi Studio website.
			</p>
			<table style="width: 100%%; border-collapse: collapse; margin-top: 20px;">
				<tr><td style="padding: 8px; font-weight: bold;">👤 Name:</td><td style="padding: 8px;">%s</td></tr>
				<tr style="background-color: #f0f2f5;"><td style="padding: 8px; font-weight: bold;">📧 Email:</td><td style="padding: 8px;">%s</td></tr>
				<tr><td style="padding: 8px; font-weight: bold;">💬 Discord:</td><td style="padding: 8px;">%s</td></tr>
				<tr style="background-color: #f0f2f5;"><td style="padding: 8px; font-weight: bold;">🎨 Type:</td><td style="padding: 8px;">%s</td></tr>
				<tr><td style="padding: 8px; font-weight: bold;">📝 Details:</td><td style="padding: 8px;">%s</td></tr>
			</table>
			<hr style="border: none; border-top: 1px solid #ddd; margin: 25px 0;">
			<p style="text-align: center; color: #777;">📅 Received on %s</p>
			<p style="text-align: center; font-size: 14px; color: #999;">
				– Plexdi Studio Notification System
			</p>
		</div>`,
		data.Name, data.Email, data.Discord, data.Type, data.Status, time.Now().Format("2006-01-02 15:04"),
	)

	// Prepare both payloads
	clientPayload := resendEmail{
		From:    "noreply@plexdistudio.com",
		To:      to,
		Subject: "🎨 New Commission Received",
		HTML:    clientHTML,
	}

	adminPayload := resendEmail{
		From:    "noreply@plexdistudio.com",
		To:      "plexdithanh@gmail.com",
		Subject: "🆕 New Commission Notification",
		HTML:    adminHTML,
	}

	// Helper to send email
	send := func(payload resendEmail) error {
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+apiKey)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to send email: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			return fmt.Errorf("Resend API error: %s", resp.Status)
		}
		return nil
	}

	// Send client + admin emails (parallel)
	go func() {
		if err := send(adminPayload); err != nil {
			log.Println("❌ Admin email failed:", err)
		} else {
			log.Println("✅ Admin notified")
		}
	}()

	if err := send(clientPayload); err != nil {
		return fmt.Errorf("client email failed: %w", err)
	}

	log.Printf("✅ Client + admin email sent successfully\n")
	return nil
}
