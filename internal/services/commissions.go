package services

import (
	"encoding/json"
	"fmt"
	"os"

	"bytes"
	"html/template"

	"gopkg.in/gomail.v2"
)

type Commission struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Discord  string `json:"discord"`
	Type     string `json:"type"`
	Platform string `json:"platform"`
	Details  string `json:"details"`
}

type CommissionData struct {
	Name string
	Type string
}

// ---------------------------- public variables --------------------------

var commissions []Commission
var nextID = 1

const filePath = "data/commissions.json"

// ----------------------------- make commmissions ------------------------------

func MakeCommission(name, email, ctype, details string) Commission {
	newCommission := Commission{
		ID:      nextID,
		Name:    name,
		Email:   email,
		Type:    ctype,
		Details: details,
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

func GetAllCommissions() []Commission {
	return commissions
}

// ------------------ Send -------------------------------------------
func SendCommissionEmail(to string, data CommissionData) error {
	// HTML template
	const emailHTML = `
	<div style="font-family: Arial, sans-serif; max-width: 600px; margin: auto; padding: 25px; background-color: #f8f9fa; border-radius: 10px;">
		<h2 style="color: #2563eb; text-align: center;">🎨 New Commission Received</h2>
		<p style="font-size: 16px; color: #333;">
			Hey <b>{{.Name}}</b>,<br><br>
			We’ve received your commission request for a <b>{{.Type}}</b>.<br>
			Our team will review it and get back to you shortly with details and payment options.
		</p>
		<hr style="border: none; border-top: 1px solid #ddd; margin: 25px 0;">
		<p style="text-align: center; font-size: 14px; color: #555;">
			– The Plexdi Studio Team<br>
			<a href="https://plexdi.studio" style="color: #2563eb; text-decoration: none;">plexdi.studio</a>
		</p>
	</div>`

	// Parse template
	tmpl, err := template.New("email").Parse(emailHTML)
	if err != nil {
		return err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	// Send the email
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAIL_USER"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", "🎨 New Commission Received")
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),
		587,
		os.Getenv("EMAIL_USER"),
		os.Getenv("EMAIL_PASS"),
	)

	return d.DialAndSend(m)
}
