package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand/v2"
	"slices"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamotypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

type Request struct {
	CityID int `json:"city_id"`
}

func UpdateRandomWeather(ctx context.Context, event json.RawMessage) error {
	req := &Request{}
	cityIds := []int{1, 13, 23, 27, 40}
	cityNames := map[int]string{1: "札幌", 13: "東京", 23: "名古屋", 27: "大阪", 40: "博多"}
	if err := json.Unmarshal(event, req); err != nil {
		log.Println("error unmarshal request", err)
		return err
	}
	if !slices.Contains(cityIds, req.CityID) {
		index := rand.IntN(len(cityIds))
		req.CityID = cityIds[index]
	}

	w := types.Weather{
		CityID:   req.CityID,
		CityName: cityNames[req.CityID],
	}
	weatherProb := rand.IntN(11) * 10
	switch {
	case weatherProb > 70:
		w.WeatherID = 12
		w.WeatherName = "雨"
	case weatherProb > 20:
		w.WeatherID = 4
		w.WeatherName = "くもり"
	default:
		w.WeatherID = 2
		w.WeatherName = "晴れ"
	}

	update := expression.Set(expression.Name("weather_id"), expression.Value(w.WeatherID))
	update.Set(expression.Name("city_name"), expression.Value(w.CityName))
	update.Set(expression.Name("weather_name"), expression.Value(w.WeatherName))
	update.Set(expression.Name("rainfall_prob"), expression.Value(weatherProb))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		log.Println("error creating update query", err)
		return err
	}

	_, err = dynamoClient.UpdateItem(
		ctx, &dynamodb.UpdateItemInput{
			Key:                       w.GetPrimaryKey(),
			TableName:                 aws.String("simple-weather-data"),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			ReturnValues:              dynamotypes.ReturnValueUpdatedNew,
			UpdateExpression:          expr.Update(),
		},
	)
	if err != nil {
		log.Println("error requesting update", err)
		return err
	}
	log.Println("update weather data: ", w.CityID)

	return nil
}

func main() {
	lambda.Start(UpdateRandomWeather)
}
