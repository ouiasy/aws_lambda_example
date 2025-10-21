package dynamoutil

import (
	"context"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/ouiasy/aws_lambda_example/internal/model"
	"github.com/ouiasy/aws_lambda_example/internal/weatherutil"
)

type DynamoClient struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoClient(ctx context.Context, tableName string) *DynamoClient {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	client := dynamodb.NewFromConfig(cfg)
	return &DynamoClient{
		client, tableName,
	}
}

func (c *DynamoClient) GetWeather(ctx context.Context, id string) (*model.Weather, error) {
	item, err := c.client.GetItem(
		ctx, &dynamodb.GetItemInput{
			TableName: aws.String(c.tableName),
			Key: map[string]types.AttributeValue{
				"city_id": &types.AttributeValueMemberN{Value: id},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	weather := &model.Weather{}
	err = attributevalue.UnmarshalMap(item.Item, weather)
	if err != nil || weather == nil {
		return nil, err
	}
	return weather, nil
}

func (c *DynamoClient) PutWeather(ctx context.Context, id string, weatherId, rainProb int) error {
	cityName := weatherutil.IdToCityName(id)
	weather := weatherutil.WeatherIdToName(strconv.Itoa(weatherId))
	update := expression.
		Set(expression.Name("rainfall_prob"), expression.Value(rainProb)).
		Set(expression.Name("city_name"), expression.Value(cityName)).
		Set(expression.Name("weather_id"), expression.Value(weatherId)).
		Set(expression.Name("weather_name"), expression.Value(weather))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return err
	}
	_, err = c.client.UpdateItem(
		ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String(c.tableName),
			Key: map[string]types.AttributeValue{
				"city_id": &types.AttributeValueMemberN{Value: id},
			},
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			ReturnValues:              types.ReturnValueUpdatedNew,
			UpdateExpression:          expr.Update(),
		},
	)
	if err != nil {
		return err
	}
	return nil
}
