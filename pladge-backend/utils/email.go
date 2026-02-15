package utils

import (
	"net/smtp"
	"net/textproto"
	"pledge-backend/config"

	"github.com/jordan-wright/email"
)

func SendEmail(data []byte, dataType int) error {
	e := &email.Email{
		To:      config.Config.Email.To,      // []string{"test@example.com"},
		Cc:      config.Config.Email.Cc,      // []string{"test@example.com"},
		From:    config.Config.Email.From,    // "Jordan Wright <test@gmail.com>",
		Subject: config.Config.Email.Subject, //"Awesome Subject",
		Headers: textproto.MIMEHeader{},
	}
	if dataType == 1 {
		e.Text = data
	} else {
		e.HTML = data
	}
	return e.Send(config.Config.Email.Host+":"+config.Config.Email.Port, smtp.PlainAuth(
		"", config.Config.Email.Username, config.Config.Email.Pwd,
		config.Config.Email.Host))
}

func SendEmailWithAttach(data []byte, dataType int, filename string) error {
	e := &email.Email{
		To:      config.Config.Email.To,
		Cc:      config.Config.Email.Cc,
		From:    config.Config.Email.From,
		Subject: config.Config.Email.Subject,
		Headers: textproto.MIMEHeader{},
	}
	if dataType == 1 {
		e.Text = data
	} else {
		e.HTML = data
	}
	_, err := e.AttachFile(filename)
	if err != nil {
		return err
	}
	return e.Send(config.Config.Email.Host+":"+config.Config.Email.Port, smtp.PlainAuth("", config.Config.Email.Username, config.Config.Email.Pwd, config.Config.Email.Host))
}
