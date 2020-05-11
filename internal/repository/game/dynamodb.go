package game

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/matiasvarela/errors"
	"github.com/matiasvarela/minesweeper-API/internal/core/domain"
	"github.com/matiasvarela/minesweeper-API/pkg/apperrors"
	"github.com/matiasvarela/minesweeper-API/pkg/dynamodbiface"
)

type awsDynamoDB struct {
	client dynamodbiface.DynamoDB
}

func NewDynamoDB(client dynamodbiface.DynamoDB) *awsDynamoDB {
	return &awsDynamoDB{client: client}
}

type GameKey struct {
	ID string `json:"id"`
}

func (db *awsDynamoDB) Get(id string) (*domain.Game, error) {
	key, err := dynamodbattribute.MarshalMap(GameKey{ID: id})
	if err != nil {
		return nil, errors.New(apperrors.Internal, err, "an internal error has occurred", "failed at generating dynamo db key")
	}

	result, err := db.client.GetItem(&dynamodb.GetItemInput{Key: key, TableName: aws.String("Games")})
	if err != nil {
		return nil, errors.New(apperrors.Internal, err, "an internal error has occurred", "failed at getting item from dynamo db")
	}

	if result.Item == nil {
		return nil, nil
	}

	game := domain.Game{}
	if err := dynamodbattribute.UnmarshalMap(result.Item, &game); err != nil {
		return nil, errors.New(apperrors.Internal, err, "an internal error has occurred", "failed at unmarshalling item")
	}

	return &game, nil
}

func (db *awsDynamoDB) Save(game domain.Game) error {
	item, err := dynamodbattribute.MarshalMap(game)
	if err != nil {
		return errors.New(apperrors.Internal, err, "an internal error has occurred", "failed at creating item")
	}

	_, err = db.client.PutItem(&dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String("Games"),
	})

	if err != nil {
		return errors.New(apperrors.Internal, err, "an internal error has occurred", "failed at saving item")
	}

	return nil
}
