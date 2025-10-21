package dynamoutil

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/ouiasy/aws_lambda_example/internal/model"
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
