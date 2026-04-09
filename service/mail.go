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
	dailer := gomail.NewDialer(
		config.Env.Mailer.Host,
		config.Env.Mailer.Port,
		config.Env.Mailer.User,
		config.Env.Mailer.Pass,
	)

	// 1. Parse the HTML template file
	t, err := template.ParseFiles("templates/emails/magic-link.html")
	if err != nil {
		log.Fatal(err)
	}

	data := map[string]any{
		"Email":     email,
		"MagicLink": fmt.Sprintf("%s/key/activate?token=%s", config.Env.App.Url, token),
		"ExpiresIn": "15 minutes",
		"AppUrl":    config.Env.App.Url,
	}

	// 2. Render the template with dynamic data into a buffer
	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		log.Fatal(err)
	}

	message := gomail.NewMessage()
	message.SetHeader("From", "Smply<noreply@smply.cc>")
	message.SetHeader("To", email)
	message.SetHeader("Subject", "Your Magic Link")
	message.SetBody("text/html", body.String())

	if err := dailer.DialAndSend(message); err != nil {
		log.Printf("Failed to send email: %v", err)
	} else {
		log.Printf("Magic link email sent to %s", email)
	}
}
