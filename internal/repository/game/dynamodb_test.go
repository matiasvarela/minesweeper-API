package game_test

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/golang/mock/gomock"
	"github.com/matiasvarela/errors"
	"github.com/matiasvarela/minesweeper-API/internal/core/domain"
	"github.com/matiasvarela/minesweeper-API/internal/repository/game"
	"github.com/matiasvarela/minesweeper-API/mock"
	"github.com/matiasvarela/minesweeper-API/pkg/apperrors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type dep struct {
	client *mock.MockDynamoDB
}

func newDep(t *testing.T) dep {
	return dep{
		client: mock.NewMockDynamoDB(gomock.NewController(t)),
	}
}

func TestAwsDynamoDB_Get(t *testing.T) {
	type args struct {
		id string
	}

	type want struct {
		result *domain.Game
		err    error
	}

	tests := []struct {
		name string
		args args
		want want
		mock func(dep, args)
	}{
		{
			name: "get game successfully",
			args: args{id: "xyz"},
			want: want{result: &domain.Game{ID: "xyz"}},
			mock: func(dep dep, arg args) {
				r, _ := dynamodbattribute.MarshalMap(domain.Game{ID: "xyz"})
				dep.client.EXPECT().GetItem(gomock.Any()).Return(&dynamodb.GetItemOutput{Item: r}, nil)
			},
		},
		{
			name: "fail at getting game from dynamodb",
			args: args{id: "xyz"},
			want: want{result: nil, err: errors.New(apperrors.Internal, apperrors.Internal, "an internal error has occurred", "failed at getting item from dynamo db")},
			mock: func(dep dep, arg args) {
				dep.client.EXPECT().GetItem(gomock.Any()).Return(nil, apperrors.Internal)
			},
		},
		{
			name: "game not found",
			args: args{id: "xyz"},
			want: want{result: nil, err: nil},
			mock: func(dep dep, arg args) {
				dep.client.EXPECT().GetItem(gomock.Any()).Return(&dynamodb.GetItemOutput{Item: nil}, nil)
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			dep := newDep(t)
			repo := game.NewDynamoDB(dep.client)
			tt.mock(dep, tt.args)
			result, err := repo.Get(tt.args.id)

			assert.Equal(t, tt.want.result, result)
			if err != nil && tt.want.err != nil {
				assert.Equal(t, tt.want.err.Error(), err.Error())
			}
			assert.Equal(t, errors.Code(tt.want.err), errors.Code(err))
		})
	}
}

func TestAwsDynamoDB_Save(t *testing.T) {
	type args struct {
		game domain.Game
	}

	type want struct {
		err    error
	}

	tests := []struct {
		name string
		args args
		want want
		mock func(dep, args)
	}{
		{
			name: "save game successfully",
			args: args{game: domain.Game{ID:"xyz"}},
			want: want{err: nil},
			mock: func(dep dep, arg args) {
				dep.client.EXPECT().PutItem(gomock.Any()).Return(nil, nil)
			},
		},
		{
			name: "fail at save the game into dynamodb",
			args: args{game: domain.Game{ID:"xyz"}},
			want: want{err: errors.New(apperrors.Internal, apperrors.Internal, "an internal error has occurred", "failed at saving item")},
			mock: func(dep dep, arg args) {
				dep.client.EXPECT().PutItem(gomock.Any()).Return(nil, apperrors.Internal)
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			dep := newDep(t)
			repo := game.NewDynamoDB(dep.client)
			tt.mock(dep, tt.args)
			err := repo.Save(tt.args.game)

			if err != nil && tt.want.err != nil {
				assert.Equal(t, tt.want.err.Error(), err.Error())
			}
			assert.Equal(t, errors.Code(tt.want.err), errors.Code(err))
		})
	}
}
