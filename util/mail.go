package util

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmailOtp(toEmail string, otp int, purpose string) (err error) {
	smtpServer := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", os.Getenv("EMAIL"), os.Getenv("PASSKEY"), smtpServer)
	from := os.Getenv("EMAIL")
	to := []string{toEmail}
	subject := "OTP for Password Manager"
	body := fmt.Sprintf("OTP for %v is %v", purpose, otp)

	message := "Subject: " + subject + "\r\n" + "\r\n" + body

	return smtp.SendMail(smtpServer+":"+smtpPort, auth, from, to, []byte(message))
}
