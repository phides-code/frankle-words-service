package main

import (
	"context"
	"math/rand"
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Entity struct {
	Id   int    `json:"id" dynamodbav:"id"`
	Word string `json:"word" dynamodbav:"word"`
}

func getClient() (dynamodb.Client, error) {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())

	dbClient := *dynamodb.NewFromConfig(sdkConfig)

	return dbClient, err
}

func getRandomWord(ctx context.Context) (*Entity, error) {
	itemCount, err := getWordCount(ctx)

	if err != nil {
		return nil, err
	}

	randomId := getRandomId(itemCount)

	word, err := getWordById(ctx, randomId)
	if err != nil {
		return nil, err
	}

	return word, nil
}

func getWordCount(ctx context.Context) (int, error) {
	input := &dynamodb.DescribeTableInput{
		TableName: aws.String(TableName),
	}

	result, err := db.DescribeTable(ctx, input)
	if err != nil {
		return 0, err
	}
	return int(*result.Table.ItemCount), nil
}

func getRandomId(itemCount int) int {
	return rand.Intn(itemCount) + 1
}

func getWordById(ctx context.Context, id int) (*Entity, error) {
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

	result, err := db.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, nil
	}

	word := new(Entity)
	err = attributevalue.UnmarshalMap(result.Item, word)
	if err != nil {
		return nil, err
	}

	return word, nil
}

func checkWord(ctx context.Context, word string) (bool, error) {
	allWords, err := listEntities(ctx)
	if err != nil {
		return false, err
	}

	// find word in allWords
	index := slices.IndexFunc(allWords, func(item Entity) bool { return item.Word == strings.ToUpper(word) })

	if index == -1 {
		return false, nil
	} else {
		return true, nil
	}
}

func listEntities(ctx context.Context) ([]Entity, error) {
	entities := make([]Entity, 0)

	var token map[string]types.AttributeValue

	for {
		input := &dynamodb.ScanInput{
			TableName:         aws.String(TableName),
			ExclusiveStartKey: token,
		}

		result, err := db.Scan(ctx, input)
		if err != nil {
			return nil, err
		}

		var fetchedEntity []Entity
		err = attributevalue.UnmarshalListOfMaps(result.Items, &fetchedEntity)
		if err != nil {
			return nil, err
		}

		entities = append(entities, fetchedEntity...)
		token = result.LastEvaluatedKey
		if token == nil {
			break
		}

	}

	return entities, nil
}
