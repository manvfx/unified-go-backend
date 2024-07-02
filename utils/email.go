package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/go-gomail/gomail"
)

func GenerateVerificationCode() string {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(b)[:6]
}

func SendVerificationEmail(to string, code string, smtpHost string, smtpPort int, smtpUser string, smtpPassword string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", smtpUser)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Email Verification")
	m.SetBody("text/plain", fmt.Sprintf("Your verification code is: %s", code))

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPassword)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
