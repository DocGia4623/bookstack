package utils

import (
	"log"
	"net/smtp"
)

// Replace these with your actual Mailtrap credentials
const (
	mailtrapHost = "sandbox.smtp.mailtrap.io"
	mailtrapPort = "587"
	mailtrapUser = "42bc7aa315900e"
	mailtrapPass = "58fbf62c2699de"
)

func OrderNotificationEmail(toEmail, orderId string) error {
	from := "test@example.com"

	msg := "From: " + from + "\n" +
		"To: " + toEmail + "\n" +
		"Subject: Order Notification\n\n" +
		"You have a new order with order ID: " + orderId

	auth := smtp.PlainAuth("", mailtrapUser, mailtrapPass, mailtrapHost)

	err := smtp.SendMail(mailtrapHost+":"+mailtrapPort, auth, from, []string{toEmail}, []byte(msg))
	if err != nil {
		log.Printf("smtp error: %s", err)
		return err
	}

	log.Print("Order notification email sent successfully")
	return nil
}

func SendVerificationEmail(toEmail, username string) error {
	from := "test@example.com"

	msg := "From: " + from + "\n" +
		"To: " + toEmail + "\n" +
		"Subject: Hello " + username + "\n\n" +
		"Here is your verification link: " + "http://localhost:8080/verify?email=" + toEmail

	auth := smtp.PlainAuth("", mailtrapUser, mailtrapPass, mailtrapHost)

	err := smtp.SendMail(mailtrapHost+":"+mailtrapPort, auth, from, []string{toEmail}, []byte(msg))
	if err != nil {
		log.Printf("smtp error: %s", err)
		return err
	}

	log.Print("Verification email sent successfully")
	return nil
}
