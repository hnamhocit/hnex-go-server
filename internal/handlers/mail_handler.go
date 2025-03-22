package handlers

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"text/template"
	"time"

	gomail "gopkg.in/mail.v2"
)

type EmailData struct {
	Name  string
	Email string
	Date  time.Time
}

func SendMail(subject, to, templateName string, data any) error {
	projectRoot, err := filepath.Abs(".")
	if err != nil {
		return err
	}

	templatePath := filepath.Join(projectRoot, "internal", "templates", fmt.Sprintf("%s.html", templateName))

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return err
	}

	dialer := gomail.NewDialer(
		os.Getenv("SMTP_PROVIDER"),
		smtpPort,
		os.Getenv("SMTP_EMAIL"),
		os.Getenv("SMTP_PASSWORD"),
	)

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", os.Getenv("SMTP_EMAIL"))
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body.String())

	if err := dialer.DialAndSend(mailer); err != nil {
		log.Fatalf("Failed to send email: %v", err)
		return err
	}

	return nil
}
