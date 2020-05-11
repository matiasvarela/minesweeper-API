package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func NewProdDynamoDB() *dynamodb.DynamoDB {
	config := &aws.Config{
		Region:   aws.String("us-east-2"),
		Credentials: credentials.NewEnvCredentials(),
	}

	sess := session.Must(session.NewSession(config))

	return dynamodb.New(sess)
}

func NewLocalDynamoDB() *dynamodb.DynamoDB{
	config := &aws.Config{
		Region:   aws.String("us-west-2"),
		Endpoint: aws.String("http://localhost:8000"),
	}

	sess := session.Must(session.NewSession(config))

	svc := dynamodb.New(sess)

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String("Games"),
	}

	svc.DeleteTable(&dynamodb.DeleteTableInput{
		TableName: aws.String("Games"),
	})

	_, err := svc.CreateTable(input)
	if err != nil {
		panic(err)
	}

	return svc
}