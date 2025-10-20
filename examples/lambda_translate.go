package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/translate"
)

var (
	translateClient    *translate.Client
	sourceLanguageCode = "ja"
	targetLanguageCode = "en"
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	translateClient = translate.NewFromConfig(cfg)
}

type Order struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}

func handleRequest(ctx context.Context, event json.RawMessage) (json.RawMessage, error) {
	order := &Order{}
	if err := json.Unmarshal(event, order); err != nil {
		return nil, err
	}

	translateText, err := translateClient.TranslateText(
		ctx, &translate.TranslateTextInput{
			SourceLanguageCode: &sourceLanguageCode,
			TargetLanguageCode: &targetLanguageCode,
			Text:               &order.Name,
			Settings:           nil,
			TerminologyNames:   nil,
		},
	)
	if err != nil {
		return nil, err
	}
	fmt.Println(os.Getenv("KEY_TEST"))
	fmt.Println(*translateText.TranslatedText)
	fmt.Println("Hello, Lambda!")
	return event, nil
}

func main() {
	lambda.Start(handleRequest)
}
