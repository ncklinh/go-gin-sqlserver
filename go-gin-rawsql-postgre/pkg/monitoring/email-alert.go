package monitoring

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmailAlert(subject, body string) error {
    from :=  os.Getenv("SENDER_EMAIL")
    to := os.Getenv("RECEIVER_EMAIL")
    password := os.Getenv("GMAIL_APP_PASSWORD")

    msg := "From: " + from + "\n" +
        "To: " + to + "\n" +
        "Subject: " + subject + "\n\n" +
        body

    // Gmail SMTP server
    err := smtp.SendMail("smtp.gmail.com:587",
        smtp.PlainAuth("", from, password, "smtp.gmail.com"),
        from, []string{to}, []byte(msg))

    if err != nil {
        return fmt.Errorf("failed to send email: %w", err)
    }
    return nil
}