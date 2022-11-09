package main

import (
	"fmt"
	"net/smtp"
)

func main() {
	from := "adenipekunmoses11@gmail.com"
	to := "oyebodeamirdeen@gmail.com"
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	headers := map[string]string{
		"From":                from,
		"To":                  to,
		"Subject":             "Mail from Amirdeen",
		"MIME-Version":        "1.0",
		"Content-Type":        "text/plain; charset=utf-8;",
		"Content-Disposition": "inline",
	}

	headerMessage := ""

	for header, value := range headers {
		headerMessage += fmt.Sprintf("%s: %s\r\n", header, value)
	}

	body := "this is the message body"

	message := headerMessage + "\r\n" + body

	auth := smtp.PlainAuth("", "oyebodeamirdeen@gmail.com", "", smtpHost)

	// Sending email.
	if err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(message)); err != nil {
		panic(err)
	}
}
