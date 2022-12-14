package main

import (
	"encoding/json"
	"fmt"
	"line-bot-sandbox/linebotsdk"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
)

func printStructure(structure interface{}) {
	fmt.Printf("(%%#v) %#v\n", structure)
}

func convertLineBotRequest(requestBody string) (*linebotsdk.RequestEvent, error) {
	var event linebotsdk.RequestEvent
	err := json.Unmarshal([]byte(requestBody), &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	printStructure(request)
	err := godotenv.Load(".env")
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	channelSecret := os.Getenv("LINE_BOT_CHANNEL_SECRET")
	channelAccessToken := os.Getenv("LINE_BOT_CHANNEL_ACCESS_TOKEN")
	bot := linebotsdk.NewBot(&linebotsdk.BotConfig{
		ChannelSecret:      channelSecret,
		ChannelAccessToken: channelAccessToken,
	})
	err = bot.ReplyText(request.Body, func(event *linebotsdk.RequestEvent) string {
		return event.Message.Text
	})
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Succeeded reply message"),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
