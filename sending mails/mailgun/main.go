package main

import (
	"context"
	"fmt"
	"log"
	"time"

	mailgun "github.com/mailgun/mailgun-go/v3"
)

// Your available domain names can be found here:
// (https://app.mailgun.com/app/domains)
var mailgunEmailDomain = "domain" // e.g. mg.yourcompany.com

// You can find the Private API Key in your Account Menu, under "Settings":
// (https://app.mailgun.com/app/account/security)
var mailgunPrivateAPIKey = "private-key"

func sendEmail(body string, recipient string) error {
	// Create an instance of the Mailgun Client
	mg := mailgun.NewMailgun(mailgunEmailDomain, mailgunPrivateAPIKey)

  // TODO: set sender to the right e-mail address
	sender := "amirdeen@" + mailgunEmailDomain
	subject := "Presstige daily tasks on " + time.Now().Format("2006-01-02")

	// The message object allows you to add attachments and Bcc recipients
	message := mg.NewMessage(sender, subject, body, recipient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	// Send the message	with a 10 second timeout
	_, _, err := mg.Send(ctx, message)

	return err
}




func main() {

	if err := sendEmail("test mail", "receiver@domain.com"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("sent")
}
