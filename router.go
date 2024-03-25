package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type ResponseStructure struct {
	Data         interface{} `json:"data"`
	ErrorMessage *string     `json:"errorMessage"` // can be string or nil
}

var headers = map[string]string{
	"Access-Control-Allow-Origin":  OriginURL,
	"Access-Control-Allow-Headers": "Content-Type",
}

func router(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return processGet(ctx, req)
	case "OPTIONS":
		return processOptions()
	default:
		return clientError(http.StatusMethodNotAllowed)
	}
}

func processOptions() (events.APIGatewayProxyResponse, error) {
	additionalHeaders := map[string]string{
		"Access-Control-Allow-Methods": "OPTIONS, GET",
		"Access-Control-Max-Age":       "3600",
	}
	mergedHeaders := mergeHeaders(headers, additionalHeaders)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    mergedHeaders,
	}, nil
}

func processGet(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	word, wordParam := req.PathParameters["word"]

	log.Println("req.PathParameters:", req.PathParameters)
	if wordParam {
		log.Println("word:", word)
	} else {
		return processGetRandom(ctx)
	}
	return clientError(400)
}

func processGetRandom(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	entity, err := getRandomEntity(ctx)
	if err != nil {
		return serverError(err)
	}

	response := ResponseStructure{
		Data:         entity,
		ErrorMessage: nil,
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		return serverError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseJson),
		Headers:    headers,
	}, nil
}

// func processGetEntityById(ctx context.Context, id string) (events.APIGatewayProxyResponse, error) {
// 	log.Printf("Received GET entity request with id = %s", id)

// 	entity, err := getEntity(ctx, id)
// 	if err != nil {
// 		return serverError(err)
// 	}

// 	if entity == nil {
// 		return clientError(http.StatusNotFound)
// 	}

// 	response := ResponseStructure{
// 		Data:         entity,
// 		ErrorMessage: nil,
// 	}

// 	responseJson, err := json.Marshal(response)
// 	if err != nil {
// 		return serverError(err)
// 	}
// 	log.Printf("Successfully fetched entity %s", response.Data)

// 	return events.APIGatewayProxyResponse{
// 		StatusCode: http.StatusOK,
// 		Body:       string(responseJson),
// 		Headers:    headers,
// 	}, nil
// }
