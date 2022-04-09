package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shomali11/slacker"
)

func printEvents(eventsChannel <-chan *slacker.CommandEvent) {
	for event := range eventsChannel {
		fmt.Printf("timestamp: %v\n", event.Timestamp)
		fmt.Printf("command: %v\n", event.Command)
		fmt.Printf("parameters: %v\n", event.Parameters)
		fmt.Printf("event: %v\n", event.Event)
	}
}

func main() {
	botToken := os.Getenv("slot-bot-token")
	appToken := os.Getenv("slot-socket-token")
	
	bot := slacker.NewClient(botToken, appToken)

	go printEvents(bot.CommandEvents())

	bot.Command("ping", &slacker.CommandDefinition{
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			response.Reply("pong")
		},
	})

	context, cancle := context.WithCancel(context.TODO())
	defer cancle()
	if err := bot.Listen(context); err != nil {
		log.Fatal("Error listening to bot:", err)
	}
}
