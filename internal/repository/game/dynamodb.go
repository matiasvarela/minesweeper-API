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
	tableName string
	client    dynamodbiface.DynamoDB
}

func NewDynamoDB(tableName string, client dynamodbiface.DynamoDB) *awsDynamoDB {
	return &awsDynamoDB{client: client, tableName: tableName}
}

type GameKey struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

func (db *awsDynamoDB) Get(userID string, gameID string) (*domain.Game, error) {
	key, err := dynamodbattribute.MarshalMap(GameKey{ID: gameID, UserID: userID})
	if err != nil {
		return nil, errors.New(apperrors.Internal, err, "an internal error has occurred", "failed at generating dynamo db key")
	}

	result, err := db.client.GetItem(&dynamodb.GetItemInput{Key: key, TableName: aws.String(db.tableName)})
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

func (db *awsDynamoDB) GetAll(userID string) ([]domain.Game, error) {
	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String(db.tableName),
		KeyConditions: map[string]*dynamodb.Condition{
			"user_id": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(userID),
					},
				},
			},
		},
	}

	resp, err := db.client.Query(queryInput)
	if err != nil {
		return nil, err
	}

	games := []domain.Game{}
	var game domain.Game

	for _, item := range resp.Items {
		game = domain.Game{}
		if err := dynamodbattribute.UnmarshalMap(item, &game); err != nil {
			return nil, errors.New(apperrors.Internal, err, "an internal error has occurred", "failed at unmarshalling item")
		}

		games = append(games, game)
	}

	return games, nil
}

func (db *awsDynamoDB) Save(game domain.Game) error {
	item, err := dynamodbattribute.MarshalMap(game)
	if err != nil {
		return errors.New(apperrors.Internal, err, "an internal error has occurred", "failed at creating item")
	}

	_, err = db.client.PutItem(&dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(db.tableName),
	})

	if err != nil {
		return errors.New(apperrors.Internal, err, "an internal error has occurred", "failed at saving item")
	}

	return nil
}
