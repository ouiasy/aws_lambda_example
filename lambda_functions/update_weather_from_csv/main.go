package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	dynamoClient *dynamodb.Client
	s3Client     *s3.Client
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	dynamoClient = dynamodb.NewFromConfig(cfg)
	s3Client = s3.NewFromConfig(cfg)
}

func ReadCsvAndUpdateWeather(ctx context.Context, event events.S3Event) error {
	for _, record := range event.Records {
		object, err := s3Client.GetObject(
			ctx, &s3.GetObjectInput{
				Bucket: aws.String(record.S3.Bucket.Name),
				Key:    aws.String(record.S3.Object.Key),
			},
		)
		if err != nil {
			return fmt.Errorf("error reading object from S3: %w", err)
		}

		csvReader := csv.NewReader(object.Body)
		for {
			line, err := csvReader.Read()
			if err != nil {
				break
			}
			weatherId, _ := strconv.Atoi(line[2])
			rainfallProb, _ := strconv.Atoi(line[4])

			update := expression.Set(expression.Name("city_name"), expression.Value(line[1]))
			update.Set(expression.Name("weather_id"), expression.Value(weatherId))
			update.Set(expression.Name("weather_name"), expression.Value(line[3]))
			update.Set(expression.Name("rainfall_prob"), expression.Value(rainfallProb))
			expr, err := expression.NewBuilder().WithUpdate(update).Build()
			if err != nil {
				log.Println("error creating update query", err)
				return fmt.Errorf("error creating update query: %w", err)
			}

			cityId, _ := strconv.Atoi(line[0])
			cityIdAttVal, err := attributevalue.Marshal(cityId)
			_, err = dynamoClient.UpdateItem(
				ctx, &dynamodb.UpdateItemInput{
					Key:                       map[string]types.AttributeValue{"city_id": cityIdAttVal},
					TableName:                 aws.String("simple-weather-data"),
					ExpressionAttributeNames:  expr.Names(),
					ExpressionAttributeValues: expr.Values(),
					ReturnValues:              types.ReturnValueUpdatedNew,
					UpdateExpression:          expr.Update(),
				},
			)
			if err != nil {
				log.Println("error requesting update", err)
				return err
			}
		}

		object.Body.Close()
	}

	return nil
}

func main() {
	lambda.Start(ReadCsvAndUpdateWeather)
}
