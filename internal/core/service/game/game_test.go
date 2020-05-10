package game_test

import (
	"github.com/golang/mock/gomock"
	"github.com/matiasvarela/errors"
	"github.com/matiasvarela/minesweeper-API/internal/core/domain"
	"github.com/matiasvarela/minesweeper-API/internal/core/service/game"
	"github.com/matiasvarela/minesweeper-API/mock"
	"github.com/matiasvarela/minesweeper-API/pkg/apperrors"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	e = domain.ElementEmpty
	X = domain.ElementEmptyMarked
	E = domain.ElementEmptyRevealed

	b = domain.ElementBomb
	Y = domain.ElementBombMarked
	B = domain.ElementBombRevealed
)

type dep struct {
	rnd        *mock.MockRandom
	repository *mock.MockGameRepository
}

func newDep(t *testing.T) dep {
	return dep{
		rnd:        mock.NewMockRandom(gomock.NewController(t)),
		repository: mock.NewMockGameRepository(gomock.NewController(t)),
	}
}

func TestService_Get(t *testing.T) {
	game1 := EasyMockGame("xyz", domain.GameStateNew)

	type args struct {
		id string
	}
	type want struct {
		result domain.Game
		err    error
	}

	tests := []struct {
		name string
		args args
		want want
		mock func(dep dep, args args)
	}{
		{
			name: "get game successfully",
			args: args{id: game1.ID},
			want: want{result: game1},
			mock: func(dep dep, args args) {
				dep.repository.EXPECT().Get(args.id).Return(&game1, nil)
			},
		},
		{
			name: "game not found",
			args: args{id: "any"},
			want: want{result: domain.Game{}, err: errors.New(apperrors.NotFound, nil, "game has not been found", "")},
			mock: func(dep dep, args args) {
				dep.repository.EXPECT().Get(args.id).Return(nil, nil)
			},
		},
		{
			name: "fail at get from repository",
			args: args{id: "any"},
			want: want{result: domain.Game{}, err: errors.New(apperrors.Internal, apperrors.Internal, "an internal error has occurred", "failed at getting game from repository")},
			mock: func(dep dep, args args) {
				dep.repository.EXPECT().Get(args.id).Return(nil, apperrors.Internal)
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			dep := newDep(t)
			service := game.NewService(dep.rnd, dep.repository)
			tt.mock(dep, tt.args)
			result, err := service.Get(tt.args.id)

			assert.Equal(t, tt.want.result, result)
			if err != nil && tt.want.err != nil {
				assert.Equal(t, tt.want.err.Error(), err.Error())
			}
			assert.Equal(t, errors.Code(tt.want.err), errors.Code(err))
		})
	}
}

func TestService_Create(t *testing.T) {
	game1 := EasyMockGame("xyz", domain.GameStateNew)

	type args struct {
		settings domain.GameSettings
	}
	type want struct {
		result domain.Game
		err    error
	}

	tests := []struct {
		name string
		args args
		want want
		mock func(dep dep, args args)
	}{
		{
			name: "create game successfully",
			args: args{settings: game1.Settings},
			want: want{result: game1},
			mock: func(dep dep, args args) {
				dep.rnd.EXPECT().GenerateID().Return("xyz")
				dep.repository.EXPECT().Save(game1).Return(nil)
			},
		},
		{
			name: "fail at save in repository",
			args: args{settings: game1.Settings},
			want: want{result: domain.Game{}, err: errors.New(apperrors.Internal, apperrors.Internal, "an internal error has occurred", "failed at saving game into repository")},
			mock: func(dep dep, args args) {
				dep.rnd.EXPECT().GenerateID().Return("xyz")
				dep.repository.EXPECT().Save(game1).Return(apperrors.Internal)
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			dep := newDep(t)
			service := game.NewService(dep.rnd, dep.repository)
			tt.mock(dep, tt.args)
			result, err := service.Create(tt.args.settings)

			assert.Equal(t, tt.want.result, result)
			if err != nil && tt.want.err != nil {
				assert.Equal(t, tt.want.err.Error(), err.Error())
			}
			assert.Equal(t, errors.Code(tt.want.err), errors.Code(err))
		})
	}
}

func TestService_MarkSquare(t *testing.T) {
	game1 := EasyMockGame("xyz", domain.GameStateNew)
	game1.Board.Set(domain.NewPosition(1, 1), domain.ElementEmpty)
	gameResult1 := EasyMockGame("xyz", domain.GameStateNew)
	gameResult1.Board.Set(domain.NewPosition(1, 1), domain.ElementEmptyMarked)

	game2 := EasyMockGame("xyz", domain.GameStateNew)
	game2.Board.Set(domain.NewPosition(1, 1), domain.ElementEmptyMarked)
	gameResult2 := EasyMockGame("xyz", domain.GameStateNew)
	gameResult2.Board.Set(domain.NewPosition(1, 1), domain.ElementEmpty)

	game3 := EasyMockGame("xyz", domain.GameStateNew)
	game3.Board.Set(domain.NewPosition(1, 1), domain.ElementBomb)
	gameResult3 := EasyMockGame("xyz", domain.GameStateNew)
	gameResult3.Board.Set(domain.NewPosition(1, 1), domain.ElementBombMarked)

	game4 := EasyMockGame("xyz", domain.GameStateNew)
	game4.Board.Set(domain.NewPosition(1, 1), domain.ElementBombMarked)
	gameResult4 := EasyMockGame("xyz", domain.GameStateNew)
	gameResult4.Board.Set(domain.NewPosition(1, 1), domain.ElementBomb)

	game5 := EasyMockGame("xyz", domain.GameStateNew)
	game5.Board.Set(domain.NewPosition(1, 1), domain.ElementEmptyRevealed)
	gameResult5 := EasyMockGame("xyz", domain.GameStateNew)
	gameResult5.Board.Set(domain.NewPosition(1, 1), domain.ElementEmptyRevealed)

	type args struct {
		id     string
		row    int
		column int
	}
	type want struct {
		result domain.Game
		err    error
	}

	tests := []struct {
		name string
		args args
		want want
		mock func(dep dep, args args)
	}{
		{
			name: "mark empty cell",
			args: args{id: "xyz", row: 1, column: 1},
			want: want{result: gameResult1},
			mock: func(dep dep, args args) {
				dep.repository.EXPECT().Get(args.id).Return(&game1, nil)
				dep.repository.EXPECT().Save(gameResult1).Return(nil)
			},
		},
		{
			name: "unmark empty cell",
			args: args{id: "xyz", row: 1, column: 1},
			want: want{result: gameResult2},
			mock: func(dep dep, args args) {
				dep.repository.EXPECT().Get(args.id).Return(&game2, nil)
				dep.repository.EXPECT().Save(gameResult2).Return(nil)
			},
		},
		{
			name: "mark cell with bomb",
			args: args{id: "xyz", row: 1, column: 1},
			want: want{result: gameResult3},
			mock: func(dep dep, args args) {
				dep.repository.EXPECT().Get(args.id).Return(&game3, nil)
				dep.repository.EXPECT().Save(gameResult3).Return(nil)
			},
		},
		{
			name: "unmark cell with bomb",
			args: args{id: "xyz", row: 1, column: 1},
			want: want{result: gameResult4},
			mock: func(dep dep, args args) {
				dep.repository.EXPECT().Get(args.id).Return(&game4, nil)
				dep.repository.EXPECT().Save(gameResult4).Return(nil)
			},
		},
		{
			name: "mark revealed cell",
			args: args{id: "xyz", row: 1, column: 1},
			want: want{result: gameResult5},
			mock: func(dep dep, args args) {
				dep.repository.EXPECT().Get(args.id).Return(&game5, nil)
				dep.repository.EXPECT().Save(gameResult5).Return(nil)
			},
		},
		{
			name: "game not found",
			args: args{id: "xyz", row: 1, column: 1},
			want: want{result: domain.Game{}, err: errors.New(apperrors.NotFound, nil, "game has not been found", "")},
			mock: func(dep dep, args args) {
				dep.repository.EXPECT().Get(args.id).Return(nil, nil)
			},
		},
		{
			name: "fail at save into repository",
			args: args{id: "xyz", row: 1, column: 1},
			want: want{result: domain.Game{}, err: errors.New(apperrors.Internal, apperrors.Internal, "an internal error has occurred", "failed at saving game into repository")},
			mock: func(dep dep, args args) {
				g := EasyMockGame(args.id, domain.GameStateNew)
				dep.repository.EXPECT().Get(args.id).Return(&g, nil)
				dep.repository.EXPECT().Save(g).Return(apperrors.Internal)
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			dep := newDep(t)
			service := game.NewService(dep.rnd, dep.repository)
			tt.mock(dep, tt.args)
			result, err := service.MarkSquare(tt.args.id, tt.args.row, tt.args.column)

			assert.Equal(t, tt.want.result, result)
			if err != nil && tt.want.err != nil {
				assert.Equal(t, tt.want.err.Error(), err.Error())
			}
			assert.Equal(t, errors.Code(tt.want.err), errors.Code(err))
		})
	}
}

func EasyMockGame(id string, state string) domain.Game {
	game := domain.Game{
		ID:       id,
		Board:    domain.NewEmptyBoard(6, 6),
		Settings: domain.GameSettings{Rows: 6, Columns: 6, BombsNumber: 10},
		State:    state,
	}

	return game
}
