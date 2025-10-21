package handlers

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/ouiasy/aws_lambda_example/internal/dynamoutil"
	"github.com/ouiasy/aws_lambda_example/internal/httputil"
)

func GetCityWeather(ctx context.Context, event events.APIGatewayV2HTTPRequest) (
	events.APIGatewayV2HTTPResponse, error,
) {
	id, ok := event.PathParameters["id"]
	if !ok {
		return httputil.ErrResponse(http.StatusBadRequest, "id is required"), nil
	}

	dynamoClient := dynamoutil.NewDynamoClient(ctx, "simple-weather-data")
	weather, err := dynamoClient.GetWeather(ctx, id)
	if err != nil {
		return httputil.ErrResponse(http.StatusInternalServerError, "requested weather not found"), err
	}

	return httputil.Response(http.StatusOK, weather), nil
}
