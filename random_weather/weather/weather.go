package weather

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Weather struct {
	CityID              int    `dynamodbav:"city_id"`
	CityName            string `dynamodbav:"city_name"`
	RainfallProbability int    `dynamodbav:"rainfall_prob"`
	WeatherID           int    `dynamodbav:"weather_id"`
	WeatherName         string `dynamodbav:"weather_name"`
}

func (w *Weather) GetPrimaryKey() map[string]types.AttributeValue {
	cityId, err := attributevalue.Marshal(w.CityID)
	if err != nil {
		panic(err)
	}
	return map[string]types.AttributeValue{
		"city_id": cityId,
	}
}
