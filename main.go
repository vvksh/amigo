package amigo

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/slack-go/slack"
)

// CallHttpGetEndpoint calls a given apiEndpoint and deserializes the response to
// user provided responseObject
func CallHttpGetEndpoint(apiEndpoint string, responseObject interface{}) error {
	resp, err := http.Get(apiEndpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return readErr
	}
	return json.Unmarshal(body, responseObject)
}

// SendSlackNotification sends message to specified channel using webhookUrl
// Prereq:
// - Install "Incoming webhooks" app to your slack workplace
// - Provide the webhook url as environment variable
func SendSlackNotification(msg string, channel string) error {
	webhookURL, exists := os.LookupEnv("SLACK_WEBHOOK")
	if !exists {
		log.Panicf("Environment variable SLACK_WEBHOOK not found\n")
	}
	webHookMessage := slack.WebhookMessage{}
	webHookMessage.Text = msg
	webHookMessage.Channel = channel
	return slack.PostWebhook(webhookURL, &webHookMessage)
}
