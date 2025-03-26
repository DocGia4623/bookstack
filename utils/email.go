package utils

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

// sendVerificationEmail gửi email xác nhận cho người dùng
func SendVerificationEmail(toEmail, username string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "anhkvse182347@fpt.edu.vn") // Thay bằng email của bạn
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Verify Your Account")
	m.SetBody("text/html", fmt.Sprintf(`
		<h2>Hello, %s!</h2>
		<p>Thank you for registering. Please click the link below to verify your email:</p>
		<a href="http://localhost:8080/verify?email=%s">Verify Email</a>
	`, username, toEmail))

	d := gomail.NewDialer("smtp.gmail.com", 587, "your-email@gmail.com", "your-app-password")

	// Gửi email
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
