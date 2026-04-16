package service

import (
	"bytes"
	"fmt"
	"log"
	"smply/config"
	"text/template"

	"gopkg.in/gomail.v2"
)

func SendAPIKeyMagicLinkEmail(email string, token string) error {
	dailer := gomail.NewDialer(
		config.Env.Mailer.Host,
		config.Env.Mailer.Port,
		config.Env.Mailer.User,
		config.Env.Mailer.Pass,
	)

	t, err := template.ParseFiles("templates/emails/magic-link.html")
	if err != nil {
		log.Printf("Failed to parse email template: %v", err)
		return err
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
		log.Printf("Error executing template: %v", err)
		return err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", "Smply<noreply@smply.cc>")
	message.SetHeader("To", email)
	message.SetHeader("Subject", "Your Magic Link")
	message.SetBody("text/html", body.String())

	if err := dailer.DialAndSend(message); err != nil {
		log.Printf("Failed to send magic link email to %s: %v", email, err)
		return err
	}

	log.Printf("Magic link email sent to %s", email)
	return nil
}

func SendAdminLoginMagicLinkEmail(email string, token string) error {
	dailer := gomail.NewDialer(
		config.Env.Mailer.Host,
		config.Env.Mailer.Port,
		config.Env.Mailer.User,
		config.Env.Mailer.Pass,
	)

	t, err := template.ParseFiles("templates/emails/magic-link.html")
	if err != nil {
		log.Printf("Failed to parse email template: %v", err)
		return err
	}

	data := map[string]any{
		"Email":     email,
		"MagicLink": fmt.Sprintf("%s/admin/auth/redirect?token=%s", config.Env.App.Url, token),
		"ExpiresIn": "15 minutes",
		"AppUrl":    config.Env.App.Url,
	}

	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		log.Printf("Error executing template: %v", err)
		return err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", "Smply<noreply@smply.cc>")
	message.SetHeader("To", email)
	message.SetHeader("Subject", "Your Admin Magic Link")
	message.SetBody("text/html", body.String())

	if err := dailer.DialAndSend(message); err != nil {
		log.Printf("Failed to send admin magic link email to %s: %v", email, err)
		return err
	}

	log.Printf("Admin magic link email sent to %s", email)
	return nil
}
