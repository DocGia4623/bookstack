package utils

import (
	"log"
	"net/smtp"
)

// sendVerificationEmail gửi email xác nhận cho người dùng
func SendVerificationEmail(toEmail, username string) error {
	from := "vietanhbestzed@gmail.com"
	pass := "oolp dail hsdf nicr"

	msg := "From: " + from + "\n" +
		"To: " + toEmail + "\n" +
		"Subject: Hello " + username + "\n\n" +
		"Here is your verification link: " + "http://localhost:8080/verify?email=" + toEmail

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{toEmail}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return err
	}

	log.Print("Email sent successfully")
	return nil
}
