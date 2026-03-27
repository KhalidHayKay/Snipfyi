package service

import (
	"bytes"
	"fmt"
	"log"
	"smply/config"
	"text/template"

	"gopkg.in/gomail.v2"
)

func sendMagicLinkEmail(email string, token string) {
	dailer := gomail.NewDialer("sandbox.smtp.mailtrap.io", 587, "45f2a26015facd", "d96a6f0d7d6b1c")

	// 1. Parse the HTML template file
	t, err := template.ParseFiles("templates/emails/magic-link.html")
	if err != nil {
		log.Fatal(err)
	}

	data := map[string]any{
		"Email":     "user@example.com",
		"MagicLink": fmt.Sprintf("%s/key/activate?token=%s", config.Env.AppUrl, token),
		"ExpiresIn": "24 hours",
		"AppUrl":    config.Env.AppUrl,
		// "RequestIP": r.RemoteAddr,
	}

	// 2. Render the template with dynamic data into a buffer
	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		log.Fatal(err)
	}

	message := gomail.NewMessage()
	message.SetHeader("From", "noreply@smply.com")
	message.SetHeader("To", email)
	message.SetHeader("Subject", "Your Magic Link")
	message.SetBody("text/html", body.String())

	if err := dailer.DialAndSend(message); err != nil {
		log.Printf("Failed to send email: %v", err)
	}
}
