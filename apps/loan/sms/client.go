package sms

import (
	"fmt"

	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

// TwilioClient struct to encapsulate Twilio's API
type TwilioClient struct {
	client     *twilio.RestClient
	fromNumber string
}

// NewTwilioClient initializes a new Twilio client
func NewTwilioClient(accountSID, authToken, fromNumber string) *TwilioClient {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})

	return &TwilioClient{
		client:     client,
		fromNumber: fromNumber,
	}
}

// SendSMS sends an SMS using Twilio API
func (t *TwilioClient) SendSMS(to, message string) error {
	params := &api.CreateMessageParams{}
	params.SetBody(message)
	params.SetFrom(t.fromNumber)
	params.SetTo(to)

	resp, err := t.client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %v", err)
	}

	// Logging response
	if resp.Sid != nil {
		fmt.Printf("SMS sent successfully! SID: %s\n", *resp.Sid)
	}

	return nil
}
