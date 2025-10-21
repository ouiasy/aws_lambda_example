package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/ouiasy/aws_lambda_example/internal/dynamoutil"
	"github.com/ouiasy/aws_lambda_example/internal/httputil"
)

func PutCityWeather(ctx context.Context, event events.APIGatewayV2HTTPRequest) (
	events.APIGatewayV2HTTPResponse, error,
) {
	id, ok := event.PathParameters["id"]
	if !ok {
		fmt.Errorf("id is required")
		return httputil.ErrResponse(http.StatusBadRequest, "id is required"), nil
	}

	req := &PutCityWeatherRequest{}
	err := json.Unmarshal([]byte(event.Body), req)
	if err != nil || req.WeatherId == nil || req.RainfallProbability == nil {
		fmt.Errorf("invalid request: %v", err)
		return httputil.ErrResponse(
			http.StatusBadRequest, "invalid request",
		), errors.New("invalid request: " + err.Error())
	}

	dynamoClient := dynamoutil.NewDynamoClient(ctx, "simple-weather-data")
	err = dynamoClient.PutWeather(ctx, id, *req.WeatherId, *req.RainfallProbability)
	if err != nil {
		fmt.Errorf("error updating weather: %v", err)
		return httputil.ErrResponse(http.StatusInternalServerError, "error updating weather"), err
	}
	return httputil.Response(http.StatusOK, nil), nil
}

type PutCityWeatherRequest struct {
	WeatherId           *int `json:"weather_id"`
	RainfallProbability *int `json:"rainfall_prob"`
}
