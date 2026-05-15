package mail

import (
	"testing"

	"smply/config"

	"gopkg.in/gomail.v2"
)

func ensureMailConfig() {
	if config.Env == nil {
		config.Env = &config.EnvType{App: config.AppConfig{Url: "http://example.com"}}
	}
	if config.Env.App.Url == "" {
		config.Env.App.Url = "http://example.com"
	}
}

func TestSendAPIKeyMagicLinkReturnsErrorWhenDialFails(t *testing.T) {
	ensureMailConfig()
	oldUrl := config.Env.App.Url
	config.Env.App.Url = "http://example.com"
	defer func() { config.Env.App.Url = oldUrl }()

	dailer := gomail.NewDialer("127.0.0.1", 1, "user", "pass")
	svc := NewService(dailer)
	if err := svc.SendAPIKeyMagicLink("user@example.com", "token"); err == nil {
		t.Fatal("expected error when dialer cannot connect")
	}
}

func TestSendAdminLoginMagicLinkReturnsErrorWhenDialFails(t *testing.T) {
	ensureMailConfig()
	oldUrl := config.Env.App.Url
	config.Env.App.Url = "http://example.com"
	defer func() { config.Env.App.Url = oldUrl }()

	dailer := gomail.NewDialer("127.0.0.1", 1, "user", "pass")
	svc := NewService(dailer)
	if err := svc.SendAdminLoginMagicLink("admin@example.com", "token"); err == nil {
		t.Fatal("expected error when dialer cannot connect")
	}
}
