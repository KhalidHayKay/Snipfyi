package service

import (
	"errors"

	"gopkg.in/gomail.v2"
)

type Mailer struct {
	Dailer *gomail.Dialer
}

func (m *Mailer) MagicLink(email, token string) error {
	m.Dailer.DialAndSend()
	return errors.New("foo")
}

func (m *Mailer) AdminLoginLink(email, token string) error {

	m.Dailer.DialAndSend()
	return errors.New("bar")
}
