package dynamodbiface

import "github.com/aws/aws-sdk-go/service/dynamodb"

//go:generate mockgen -source=dynamodbiface.go -destination=../../mock/dynamodbiface.go -package=mock

type DynamoDB interface {
	GetItem(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
	PutItem(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
}