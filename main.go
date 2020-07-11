package amigo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/kennygrant/sanitize"

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

// Strips html tags, replace common entities, and escapes <>&;'" in the result.
func Sanitize(input string) string {
	return sanitize.HTML(input)
}

func GetRHMobileStockQuoteUrl(stock string) string {
	return fmt.Sprintf("https://robinhood.com/applink/instrument/?symbol=%s", stock)
}

func GetRHWebStockQuoteUrl(stock string) string {
	return fmt.Sprintf("https://robinhood.com/stocks/%s", stock)
}
