package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/ouiasy/aws_lambda_example/internal/types"
)

var (
	dynamoClient *dynamodb.Client
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	dynamoClient = dynamodb.NewFromConfig(cfg)
}

func GetAllWeather(
	ctx context.Context, event events.APIGatewayV2HTTPRequest,
) (events.APIGatewayV2HTTPResponse, error) {
	result, err := dynamoClient.Scan(
		ctx, &dynamodb.ScanInput{
			TableName: aws.String("simple-weather-data"),
		},
	)
	if err != nil {
		log.Println(err)
		return errResponse(http.StatusInternalServerError, err.Error()), err
	}
	allWeather := &types.AllWeather{}
	err = attributevalue.UnmarshalListOfMaps(result.Items, allWeather)
	if err != nil {
		log.Println(err)
		return errResponse(http.StatusInternalServerError, err.Error()), err
	}

	return response(http.StatusOK, allWeather), nil
}

func main() {
	lambda.Start(GetAllWeather)
}

func response(code int, object interface{}) events.APIGatewayV2HTTPResponse {
	marshalled, err := json.Marshal(object)
	if err != nil {
		return errResponse(http.StatusInternalServerError, err.Error())
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: code,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:            string(marshalled),
		IsBase64Encoded: false,
	}
}

func errResponse(status int, body string) events.APIGatewayV2HTTPResponse {
	message := map[string]string{
		"message": body,
	}

	messageBytes, _ := json.Marshal(&message)

	return events.APIGatewayV2HTTPResponse{
		StatusCode: status,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(messageBytes),
	}
}
