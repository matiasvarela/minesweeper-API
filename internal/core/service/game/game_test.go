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
	"time"
)

var (
	e = domain.EmptyCellCovered
	X = domain.EmptyCellCoveredAndMarked
	E = domain.EmptyCellRevealed

	b = domain.BombCellCovered
	Y = domain.BombCellCoveredAndMarked
	B = domain.BombCellRevealed
)

type dep struct {
	rnd        *mock.MockRandom
	clock      *mock.MockClock
	repository *mock.MockGameRepository
}

func newDep(t *testing.T) dep {
	return dep{
		rnd:        mock.NewMockRandom(gomock.NewController(t)),
		clock:      mock.NewMockClock(gomock.NewController(t)),
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
			service := game.NewService(dep.rnd, dep.clock, dep.repository)
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
			service := game.NewService(dep.rnd, dep.clock, dep.repository)
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

func TestService_MarkCell(t *testing.T) {
	game1 := EasyMockGameWithElement("xyz", domain.GameStateNew, 1, 1, domain.EmptyCellCovered)
	gameResult1 := EasyMockGameWithElement("xyz", domain.GameStateNew, 1, 1, domain.EmptyCellCoveredAndMarked)

	game2 := EasyMockGameWithElement("xyz", domain.GameStateNew, 1, 1, domain.EmptyCellCoveredAndMarked)
	gameResult2 := EasyMockGameWithElement("xyz", domain.GameStateNew, 1, 1, domain.EmptyCellCovered)

	game3 := EasyMockGameWithElement("xyz", domain.GameStateNew, 1, 1, domain.BombCellCovered)
	gameResult3 := EasyMockGameWithElement("xyz", domain.GameStateNew, 1, 1, domain.BombCellCoveredAndMarked)

	game4 := EasyMockGameWithElement("xyz", domain.GameStateNew, 1, 1, domain.BombCellCoveredAndMarked)
	gameResult4 := EasyMockGameWithElement("xyz", domain.GameStateNew, 1, 1, domain.BombCellCovered)

	game5 := EasyMockGameWithElement("xyz", domain.GameStateNew, 1, 1, domain.EmptyCellRevealed)
	gameResult5 := EasyMockGameWithElement("xyz", domain.GameStateNew, 1, 1, domain.EmptyCellRevealed)

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
		{
			name: "invalid row and column params",
			args: args{id: "xyz", row: 100, column: 1},
			want: want{result: domain.Game{}, err: errors.New(apperrors.InvalidInput, apperrors.InvalidInput, "invalid row and column parameters", "")},
			mock: func(dep dep, args args) {
				g := EasyMockGame(args.id, domain.GameStateNew)
				dep.repository.EXPECT().Get(args.id).Return(&g, nil)
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			dep := newDep(t)
			service := game.NewService(dep.rnd, dep.clock, dep.repository)
			tt.mock(dep, tt.args)
			result, err := service.MarkCell(tt.args.id, tt.args.row, tt.args.column)

			assert.Equal(t, tt.want.result, result)
			if err != nil && tt.want.err != nil {
				assert.Equal(t, tt.want.err.Error(), err.Error())
			}
			assert.Equal(t, errors.Code(tt.want.err), errors.Code(err))
		})
	}
}

func TestService_RevealCell(t *testing.T) {
	mockedTime, _ := time.Parse(time.RFC3339, time.RFC3339)

	game1 := EasyMockGameWith("xyz", domain.GameStateNew, 5, domain.Board{
		{e, e, e, e, e, e},
		{e, e, e, e, e, e},
		{e, e, e, e, e, e},
		{e, e, e, e, e, e},
	}, mockedTime, time.Time{})
	gameResult1 := EasyMockGameWith("xyz", domain.GameStateOnGoing, 5, domain.Board{
		{b, e, b, e, e, e},
		{e, b, e, e, e, e},
		{e, e, E, e, b, e},
		{e, e, e, e, e, b},
	}, mockedTime, time.Time{})

	game2 := EasyMockGameWith("xyz", domain.GameStateOnGoing, 2, domain.Board{
		{e, e, e, e, e, e},
		{e, e, e, e, e, e},
		{e, e, e, b, e, e},
		{e, e, e, e, e, b},
	}, mockedTime, time.Time{})

	gameResult2 := EasyMockGameWith("xyz", domain.GameStateOnGoing, 2, domain.Board{
		{E, E, E, E, E, E},
		{E, E, E, E, E, E},
		{E, E, E, b, E, E},
		{E, E, E, e, e, b},
	}, mockedTime, time.Time{})

	game3 := EasyMockGameWith("xyz", domain.GameStateOnGoing, 2, domain.Board{
		{e, e, e, e, e, e},
		{e, e, e, e, e, e},
		{e, e, e, b, e, e},
		{e, e, e, e, e, b},
	}, mockedTime, time.Time{})

	gameResult3 := EasyMockGameWith("xyz", domain.GameStateLost, 2, domain.Board{
		{e, e, e, e, e, e},
		{e, e, e, e, e, e},
		{e, e, e, B, e, e},
		{e, e, e, e, e, b},
	}, mockedTime, mockedTime)

	game4 := EasyMockGameWith("xyz", domain.GameStateOnGoing, 2, domain.Board{
		{e, e, e, e, e, e},
		{e, e, e, e, e, e},
		{e, e, e, b, e, e},
		{e, e, e, E, E, b},
	}, mockedTime, time.Time{})

	gameResult4 := EasyMockGameWith("xyz", domain.GameStateWon, 2, domain.Board{
		{E, E, E, E, E, E},
		{E, E, E, E, E, E},
		{E, E, E, b, E, E},
		{E, E, E, E, E, b},
	}, mockedTime, mockedTime)

	game5 := EasyMockGameWithElement("xyz", domain.GameStateOnGoing, 1, 1, domain.EmptyCellCoveredAndMarked)

	game6 := EasyMockGameWith("xyz", domain.GameStateOnGoing, 2, domain.Board{
		{e, e, e, e, e, e},
		{e, e, e, e, e, e},
		{e, e, e, b, e, e},
		{e, e, e, E, E, b},
	}, mockedTime, time.Time{})

	gameResult6 := EasyMockGameWith("xyz", domain.GameStateWon, 2, domain.Board{
		{E, E, E, E, E, E},
		{E, E, E, E, E, E},
		{E, E, E, b, E, E},
		{E, E, E, E, E, b},
	}, mockedTime, mockedTime)

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
			name: "game not found",
			args: args{id: "xyz", row: 1, column: 1},
			want: want{result: domain.Game{}, err: errors.New(apperrors.NotFound, nil, "game has not been found", "")},
			mock: func(dep dep, args args) {
				dep.repository.EXPECT().Get(args.id).Return(nil, nil)
			},
		},
		{
			name: "game has already been finished - lost",
			args: args{id: "xyz", row: 1, column: 1},
			want: want{result: domain.Game{}, err: errors.New(apperrors.InvalidInput, nil, "game has already finished", "")},
			mock: func(dep dep, args args) {
				g := EasyMockGame("xyz", domain.GameStateLost)
				dep.repository.EXPECT().Get(args.id).Return(&g, nil)
			},
		},
		{
			name: "game has already been finished - won",
			args: args{id: "xyz", row: 1, column: 1},
			want: want{result: domain.Game{}, err: errors.New(apperrors.InvalidInput, nil, "game has already finished", "")},
			mock: func(dep dep, args args) {
				g := EasyMockGame("xyz", domain.GameStateWon)
				dep.repository.EXPECT().Get(args.id).Return(&g, nil)
			},
		},
		{
			name: "invalid position",
			args: args{id: "xyz", row: 100, column: 1},
			want: want{result: domain.Game{}, err: errors.New(apperrors.InvalidInput, nil, "invalid row and column parameters", "")},
			mock: func(dep dep, args args) {
				g := EasyMockGame("xyz", domain.GameStateNew)
				dep.repository.EXPECT().Get(args.id).Return(&g, nil)
			},
		},
		{
			name: "cell is marked",
			args: args{id: "xyz", row: 1, column: 1},
			want: want{result: game5},
			mock: func(dep dep, args args) {
				dep.repository.EXPECT().Get(args.id).Return(&game5, nil)
			},
		},
		{
			name: "reveal first cell successfully",
			args: args{id: "xyz", row: 2, column: 2},
			want: want{result: gameResult1},
			mock: func(dep dep, args args) {
				dep.clock.EXPECT().Now().Return(mockedTime)
				dep.rnd.EXPECT().GenerateN(game1.Settings.BombsNumber + 1).Return([]int{2, 7, 16, 14, 23, 0})
				dep.repository.EXPECT().Get(args.id).Return(&game1, nil)
				dep.repository.EXPECT().Save(gameResult1).Return(nil)
			},
		},
		{
			name: "reveal cell in cascade successfully",
			args: args{id: "xyz", row: 1, column: 1},
			want: want{result: gameResult2},
			mock: func(dep dep, args args) {
				dep.clock.EXPECT().Now().Return(mockedTime)
				dep.repository.EXPECT().Get(args.id).Return(&game2, nil)
				dep.repository.EXPECT().Save(gameResult2).Return(nil)
			},
		},
		{
			name: "reveal cell with bomb and lost game",
			args: args{id: "xyz", row: 2, column: 3},
			want: want{result: gameResult3},
			mock: func(dep dep, args args) {
				dep.clock.EXPECT().Now().Return(mockedTime)
				dep.repository.EXPECT().Get(args.id).Return(&game3, nil)
				dep.repository.EXPECT().Save(gameResult3).Return(nil)
			},
		},
		{
			name: "reveal cell and won game",
			args: args{id: "xyz", row: 1, column: 1},
			want: want{result: gameResult4},
			mock: func(dep dep, args args) {
				dep.clock.EXPECT().Now().Return(mockedTime)
				dep.repository.EXPECT().Get(args.id).Return(&game4, nil)
				dep.repository.EXPECT().Save(gameResult4).Return(nil)
			},
		},
		{
			name: "fail at save in repository",
			args: args{id: "xyz", row: 1, column: 1},
			want: want{result: domain.Game{}, err: errors.New(apperrors.Internal, apperrors.Internal, "an internal error has occurred", "failed at saving game into repository")},
			mock: func(dep dep, args args) {
				dep.clock.EXPECT().Now().Return(mockedTime)
				dep.repository.EXPECT().Get(args.id).Return(&game6, nil)
				dep.repository.EXPECT().Save(gameResult6).Return(apperrors.Internal)
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			dep := newDep(t)
			service := game.NewService(dep.rnd, dep.clock, dep.repository)
			tt.mock(dep, tt.args)
			result, err := service.RevealCell(tt.args.id, tt.args.row, tt.args.column)

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

func EasyMockGameWith(id string, state string, bombsNumber int, board domain.Board, startedAt time.Time, endedAt time.Time) domain.Game {
	game := domain.Game{
		ID:       id,
		Board:    board,
		Settings: domain.GameSettings{Rows: len(board), Columns: len(board[0]), BombsNumber: bombsNumber},
		State:    state,
	}

	return game
}

func EasyMockGameWithElement(id string, state string, row int, column int, element domain.Cell) domain.Game {
	game := domain.Game{
		ID:       id,
		Board:    domain.NewEmptyBoard(6, 6),
		Settings: domain.GameSettings{Rows: 6, Columns: 6, BombsNumber: 10},
		State:    state,
	}

	game.Board.Set(domain.NewPosition(row, column), element)

	return game
}
