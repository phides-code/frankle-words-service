package main

import (
	"context"
	"log"
	"math/rand"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Entity struct {
	// example database table structure:
	Id   int    `json:"id" dynamodbav:"id"`
	Word string `json:"word" dynamodbav:"word"`
	// adjust fields as needed
}

type NewOrUpdatedEntity struct {
	Word string `json:"description" validate:"required"`
	// adjust fields as needed
}

func getClient() (dynamodb.Client, error) {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())

	dbClient := *dynamodb.NewFromConfig(sdkConfig)

	return dbClient, err
}

func getRandomId(numberOfItems int) int {
	randomId := rand.Intn(numberOfItems) // Generates a random integer in the range [0, numberOfItems)
	randomId++                           // Adjust the index to be in the range [1, numberOfItems]

	return randomId
}

func getRandomEntity(ctx context.Context) (*Entity, error) {
	numberOfItems, err := getItemCount(ctx)
	if err != nil {
		return nil, err
	}

	randomId := getRandomId(numberOfItems)

	entity, err := getEntity(ctx, randomId)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func getItemCount(ctx context.Context) (int, error) {
	input := &dynamodb.DescribeTableInput{
		TableName: aws.String("AppnameApples"),
	}

	result, err := db.DescribeTable(ctx, input)
	if err != nil {
		return 0, err
	}

	count := result.Table.ItemCount

	return int(*count), nil
}

func getEntity(ctx context.Context, id int) (*Entity, error) {
	key, err := attributevalue.Marshal(id)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key: map[string]types.AttributeValue{
			"id": key,
		},
	}

	log.Printf("Calling DynamoDB with input: %v", input)
	result, err := db.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}
	log.Printf("Executed GetEntity DynamoDb successfully. Result: %#v", result)

	if result.Item == nil {
		return nil, nil
	}

	entity := new(Entity)
	err = attributevalue.UnmarshalMap(result.Item, entity)
	if err != nil {
		return nil, err
	}

	return entity, nil
}
