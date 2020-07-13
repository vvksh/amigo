package amigo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/sheets/v4"

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

// Appends to first sheet of the spreadsheet
func AppendToSheet(sheetsService *sheets.Service, sheetsId string, values []interface{}) error {
	writeRange := "Sheet1"

	var vr sheets.ValueRange
	vr.MajorDimension = "ROWS"
	vr.Values = append(vr.Values, values)

	_, err := sheetsService.Spreadsheets.Values.Append(sheetsId, writeRange, &vr).ValueInputOption("RAW").Do()
	sheetsService.Spreadsheets.Values.Get(sheetsId, writeRange)
	return err
}

func GetAllSheetData(sheetsService *sheets.Service, sheetsId string) ([][]interface{}, error) {
	writeRange := "Sheet1"
	valueRange, err := sheetsService.Spreadsheets.Values.Get(sheetsId, writeRange).Do()
	return valueRange.Values, err
}

func GetGmailService() (*gmail.Service, error) {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	return gmail.New(client)

}

func GetSheetsService() (*sheets.Service, error) {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	return sheets.New(client)
}
