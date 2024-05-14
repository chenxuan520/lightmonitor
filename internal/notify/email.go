package notify

import (
	"fmt"

	"github.com/chenxuan520/lightmonitor/internal/config"
	mail "github.com/wneessen/go-mail"
)

type Email struct {
	AbstractNotify
	Domain    string
	Password  string
	SendEmail string
	RecvEmail string
}

func NewEmail() *Email {
	return &Email{
		Domain:    config.GlobalConfig.NotifyWay.Email.Domain,
		Password:  config.GlobalConfig.NotifyWay.Email.Password,
		SendEmail: config.GlobalConfig.NotifyWay.Email.SendEmail,
		RecvEmail: config.GlobalConfig.NotifyWay.Email.RecvEmail,
	}
}

func (e *Email) Send(msg NotifyMsg) error {
	if e.Domain == "" {
		return fmt.Errorf("email domain is empty")
	}
	m := mail.NewMsg()
	err := m.From(e.SendEmail)
	if err != nil {
		return err
	}
	err = m.To(e.RecvEmail)
	if err != nil {
		return err
	}
	m.Subject(msg.Title)
	m.SetBodyString(mail.TypeTextPlain, msg.Content)
	c, err := mail.NewClient(config.GlobalConfig.NotifyWay.Email.Domain, mail.WithPort(25), mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(config.GlobalConfig.NotifyWay.Email.SendEmail), mail.WithPassword(config.GlobalConfig.NotifyWay.Email.Password))
	if err != nil {
		return err
	}
	if err := c.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
