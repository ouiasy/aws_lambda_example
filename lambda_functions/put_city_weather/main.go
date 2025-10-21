package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ouiasy/aws_lambda_example/internal/handlers"
)

func main() {
	lambda.Start(handlers.PutCityWeather)
}
