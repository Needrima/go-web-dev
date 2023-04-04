package main

import (
	"encoding/json"
	"fmt"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

func main() {
	from := "" // twilio phone number
	to := "" // recipient phone number
	body := "Hello world"

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: "", // Account SID
		Password: "", //Auth Token
	})

	params := &twilioApi.CreateMessageParams{
		To: &to,
		From: &from,
		Body: &body,
	}

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println(err.Error())
		return
	} 
	
	response, _ := json.Marshal(*resp)
	fmt.Println("Response: " + string(response))
}
