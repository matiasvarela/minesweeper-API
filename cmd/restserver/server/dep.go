package server

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	gameService "github.com/matiasvarela/minesweeper-API/internal/core/service/game"
	"github.com/matiasvarela/minesweeper-API/internal/dep"
	"github.com/matiasvarela/minesweeper-API/internal/handler"
	gameRepo "github.com/matiasvarela/minesweeper-API/internal/repository/game"
	"github.com/matiasvarela/minesweeper-API/pkg/clock"
	"github.com/matiasvarela/minesweeper-API/pkg/random"
	"os"
)

func initDependencies() *dep.Dep {
	d := &dep.Dep{}

	if os.Getenv("ENV") == "production" {
		d.DynamoDB = newProdDynamoDB()
	} else {
		d.DynamoDB = newLocalDynamoDB()
	}

	rnd := random.NewRandom()
	rnd.Init()

	clk := clock.New()

	d.GameRepository = gameRepo.NewDynamoDB(d.DynamoDB)
	d.GameService = gameService.NewService(rnd, clk, d.GameRepository)
	d.GameHandler = handler.NewGameHandler(d.GameService)

	return d
}

func newProdDynamoDB() *dynamodb.DynamoDB {
	config := &aws.Config{
		Region:   aws.String("us-east-2"),
		Credentials: credentials.NewEnvCredentials(),
	}

	sess := session.Must(session.NewSession(config))

	return dynamodb.New(sess)
}

func newLocalDynamoDB() *dynamodb.DynamoDB{
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