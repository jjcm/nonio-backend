package utils

import (
	"fmt"
	"net/smtp"
)

var AdminEmail string
var AdminEmailPassword string

func SendEmail(to string, subject string, body string) {
	fmt.Printf("App Email: %v\n", AdminEmail)
	fmt.Printf("App Password: %v\n", AdminEmailPassword)
	msg := "From: " + AdminEmail + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body
	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", AdminEmail, AdminEmailPassword, "smtp.gmail.com"),
		AdminEmail, []string{to}, []byte(msg))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("email sent successfully:")
	fmt.Println(msg)
}
