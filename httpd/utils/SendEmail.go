package utils

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
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

// GmailService : Gmail client for sending email
var GmailService *gmail.Service

func OAuthGmailService() {
	if GmailService != nil {
		return
	}

	config := oauth2.Config{
		ClientID:     os.Getenv("EMAIL_CLIENT_ID"),
		ClientSecret: os.Getenv("EMAIL_CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost",
	}

	token := oauth2.Token{
		AccessToken:  os.Getenv("EMAIL_ACCESS_TOKEN"),
		RefreshToken: os.Getenv("EMAIL_REFRESH_TOKEN"),
		TokenType:    "Bearer",
		Expiry:       time.Now(),
	}

	var tokenSource = config.TokenSource(context.Background(), &token)

	srv, err := gmail.NewService(context.Background(), option.WithTokenSource(tokenSource))
	if err != nil {
		log.Printf("Unable to retrieve Gmail client: %v", err)
	}

	GmailService = srv
	if GmailService != nil {
		fmt.Println("Email service is initialized \n")
	}
}

func SendEmailOAUTH2(to string, body string) (bool, error) {
	OAuthGmailService()

	var message gmail.Message

	emailTo := "To: " + to + "\r\n"
	subject := "Subject: " + "Test Email form Gmail API using OAuth" + "\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	msg := []byte(emailTo + subject + mime + "\n" + body)

	message.Raw = base64.URLEncoding.EncodeToString(msg)

	// Send the message
	_, err := GmailService.Users.Messages.Send("me", &message).Do()
	if err != nil {
		return false, err
	}
	return true, nil
}
