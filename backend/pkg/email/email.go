package email

import (
	"adv/internal/config"
	"fmt"
	"net/smtp"
)

func SendVerificationCodeEmail(email, verificationCode string, config config.Config) error {
	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: Verification Code\n\nYour verification code is: %s", config.EmailFrom, email, verificationCode)

	auth := smtp.PlainAuth("", config.EmailFrom, config.EmailPassword, config.SMTPHost)

	err := smtp.SendMail(config.SMTPHost+":"+config.SMTPPort, auth, config.EmailFrom, []string{email}, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}

func SendSpamForUser(email, text string, config config.Config) error {
	msg := []byte(text)
	auth := smtp.PlainAuth("", config.EmailFrom, config.EmailPassword, config.SMTPHost)
	err := smtp.SendMail(config.SMTPHost+":"+config.SMTPPort, auth, config.EmailFrom, []string{email}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	return nil
}
