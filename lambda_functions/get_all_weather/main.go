package main

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/ouiasy/aws_lambda_example/internal/httputil"
	"github.com/ouiasy/aws_lambda_example/internal/model"
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
		return httputil.ErrResponse(http.StatusInternalServerError, err.Error()), err
	}
	allWeather := &model.AllWeather{}
	err = attributevalue.UnmarshalListOfMaps(result.Items, allWeather)
	if err != nil {
		log.Println(err)
		return httputil.ErrResponse(http.StatusInternalServerError, err.Error()), err
	}

	return httputil.Response(http.StatusOK, allWeather), nil
}

func main() {
	lambda.Start(GetAllWeather)
}
