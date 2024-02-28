package util

import (
	"fmt"
	"halodeksik-be/app/appconfig"
	"net/smtp"
	"strings"
)

type EmailUtil interface {
	SendEmail(to []string, cc []string, subject, message string) error
}

func NewEmailUtil() EmailUtil {
	return &EmailUtilImpl{}
}

type EmailUtilImpl struct{}

func (a EmailUtilImpl) SendEmail(to []string, cc []string, subject, message string) error {
	body := "From: " + appconfig.Config.MailSender + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Cc: " + strings.Join(cc, ",") + "\n" +
		"Subject: " + subject + "\n" +
		"MIME-version: 1.0;\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\n\n" +
		message

	auth := smtp.PlainAuth("", appconfig.Config.MailAddress, appconfig.Config.MailPassword, appconfig.Config.MailSmtpHost)
	smtpAddr := fmt.Sprintf("%s:%s", appconfig.Config.MailSmtpHost, appconfig.Config.MailSmtpPort)
	err := smtp.SendMail(smtpAddr, auth, appconfig.Config.MailAddress, append(to, cc...), []byte(body))
	if err != nil {
		return err
	}

	return nil

}
