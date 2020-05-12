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
	type args struct {
		userID string
		gameID string
	}
	type want struct {
		result domain.Game
		err    error
	}

	tests := []struct {
		name string
		args args
		want want
		mock func(dep, args, want)
	}{
		{
			name: "get game successfully",
			args: args{userID: "111", gameID: "xyz"},
			want: want{result: MockGame("111", "xyz", "")},
			mock: func(dep dep, args args, want want) {
				dep.repository.EXPECT().Get(args.userID, args.gameID).Return(&want.result, nil)
			},
		},
		{
			name: "game not found",
			args: args{gameID: "any"},
			want: want{result: domain.Game{}, err: errors.New(apperrors.NotFound, nil, "game has not been found", "")},
			mock: func(dep dep, args args, want want) {
				dep.repository.EXPECT().Get(args.userID, args.gameID).Return(nil, nil)
			},
		},
		{
			name: "fail at get from repository",
			args: args{gameID: "any"},
			want: want{result: domain.Game{}, err: errors.New(apperrors.Internal, apperrors.Internal, "an internal error has occurred", "failed at getting game from repository")},
			mock: func(dep dep, args args, want want) {
				dep.repository.EXPECT().Get(args.userID, args.gameID).Return(nil, apperrors.Internal)
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			dep := newDep(t)
			service := game.NewService(dep.rnd, dep.clock, dep.repository)
			tt.mock(dep, tt.args, tt.want)
			result, err := service.Get(tt.args.userID, tt.args.gameID)

			assert.Equal(t, tt.want.result, result)
			if err != nil && tt.want.err != nil {
				assert.Equal(t, tt.want.err.Error(), err.Error())
			}
			assert.Equal(t, errors.Code(tt.want.err), errors.Code(err))
		})
	}
}

func TestService_GetAll(t *testing.T) {
	type args struct {
		userID string
	}
	type want struct {
		result []domain.Game
		err    error
	}

	tests := []struct {
		name string
		args args
		want want
		mock func(dep, args, want)
	}{
		{
			name: "get games successfully",
			args: args{userID: "111"},
			want: want{result: []domain.Game{MockGame("111", "xyz", ""), MockGame("111", "xyz", "")}},
			mock: func(dep dep, args args, want want) {
				dep.repository.EXPECT().GetAll(args.userID).Return(want.result, nil)
			},
		},
		{
			name: "fail at get from repository",
			args: args{userID: "111"},
			want: want{result: nil, err: errors.New(apperrors.Internal, apperrors.Internal, "an internal error has occurred", "failed at getting game from repository")},
			mock: func(dep dep, args args, want want) {
				dep.repository.EXPECT().GetAll(args.userID).Return(nil, apperrors.Internal)
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			dep := newDep(t)
			service := game.NewService(dep.rnd, dep.clock, dep.repository)
			tt.mock(dep, tt.args, tt.want)
			result, err := service.GetAll(tt.args.userID)

			assert.Equal(t, tt.want.result, result)
			if err != nil && tt.want.err != nil {
				assert.Equal(t, tt.want.err.Error(), err.Error())
			}
			assert.Equal(t, errors.Code(tt.want.err), errors.Code(err))
		})
	}
}

func TestService_MarkCell(t *testing.T) {
	type args struct {
		userID string
		gameID string
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
		mock func(dep, args, want)
	}{
		{
			name: "mark - empty covered cell",
			args: args{userID: "111", gameID: "xyz", row: 1, column: 1},
			want: want{result: MockGameWithCell("111", "xyz", "", 1, 1, domain.EmptyCellCoveredAndMarked)},
			mock: func(dep dep, args args, want want) {
				game := MockGameWithCell("111", "xyz", "", 1, 1, domain.EmptyCellCovered)

				dep.repository.EXPECT().Get(args.userID, args.gameID).Return(&game, nil)
				dep.repository.EXPECT().Save(want.result).Return(nil)
			},
		},
		{
			name: "unmark - marked empty covered cell",
			args: args{userID: "111", gameID: "xyz", row: 1, column: 1},
			want: want{result: MockGameWithCell("111", "xyz", "", 1, 1, domain.EmptyCellCovered)},
			mock: func(dep dep, args args, want want) {
				game := MockGameWithCell("111", "xyz", "", 1, 1, domain.EmptyCellCoveredAndMarked)

				dep.repository.EXPECT().Get(args.userID, args.gameID).Return(&game, nil)
				dep.repository.EXPECT().Save(want.result).Return(nil)
			},
		},
		{
			name: "mark - cell bomb covered",
			args: args{userID: "111", gameID: "xyz", row: 1, column: 1},
			want: want{result: MockGameWithCell("111", "xyz", "", 1, 1, domain.BombCellCoveredAndMarked)},
			mock: func(dep dep, args args, want want) {
				game := MockGameWithCell("111", "xyz", "", 1, 1, domain.BombCellCovered)

				dep.repository.EXPECT().Get(args.userID, args.gameID).Return(&game, nil)
				dep.repository.EXPECT().Save(want.result).Return(nil)
			},
		},
		{
			name: "unmark - marked cell bomb covered",
			args: args{userID: "111", gameID: "xyz", row: 1, column: 1},
			want: want{result: MockGameWithCell("111", "xyz", "", 1, 1, domain.BombCellCovered)},
			mock: func(dep dep, args args, want want) {
				game := MockGameWithCell("111", "xyz", "", 1, 1, domain.BombCellCoveredAndMarked)

				dep.repository.EXPECT().Get(args.userID, args.gameID).Return(&game, nil)
				dep.repository.EXPECT().Save(want.result).Return(nil)
			},
		},
		{
			name: "mark - revealed empty cell",
			args: args{userID: "111", gameID: "xyz", row: 1, column: 1},
			want: want{result: MockGameWithCell("111", "xyz", "", 1, 1, domain.BombCellRevealed)},
			mock: func(dep dep, args args, want want) {
				game := MockGameWithCell("111", "xyz", "", 1, 1, domain.BombCellRevealed)

				dep.repository.EXPECT().Get(args.userID, args.gameID).Return(&game, nil)
				dep.repository.EXPECT().Save(want.result).Return(nil)
			},
		},
		{
			name: "game not found",
			args: args{userID: "111", gameID: "xyz", row: 1, column: 1},
			want: want{result: domain.Game{}, err: errors.New(apperrors.NotFound, nil, "game has not been found", "")},
			mock: func(dep dep, args args, want want) {
				dep.repository.EXPECT().Get(args.userID, args.gameID).Return(nil, nil)
			},
		},
		{
			name: "fail at save into repository",
			args: args{userID: "111", gameID: "xyz", row: 1, column: 1},
			want: want{result: domain.Game{}, err: errors.New(apperrors.Internal, apperrors.Internal, "an internal error has occurred", "failed at saving game into repository")},
			mock: func(dep dep, args args, want want) {
				game := MockGame(args.userID, args.gameID, "")
				dep.repository.EXPECT().Get(args.userID, args.gameID).Return(&game, nil)
				dep.repository.EXPECT().Save(game).Return(apperrors.Internal)
			},
		},
		{
			name: "invalid row and column params",
			args: args{userID: "111", gameID: "xyz", row: -100, column: 1},
			want: want{result: domain.Game{}, err: errors.New(apperrors.InvalidInput, apperrors.InvalidInput, "invalid row and column parameters", "")},
			mock: func(dep dep, args args, want want) {
				game := MockGame(args.userID, args.gameID, "")
				dep.repository.EXPECT().Get(args.userID, args.gameID).Return(&game, nil)
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			dep := newDep(t)
			service := game.NewService(dep.rnd, dep.clock, dep.repository)
			tt.mock(dep, tt.args, tt.want)
			result, err := service.MarkCell(tt.args.userID, tt.args.gameID, tt.args.row, tt.args.column)

			assert.Equal(t, tt.want.result, result)
			if err != nil && tt.want.err != nil {
				assert.Equal(t, tt.want.err.Error(), err.Error())
			}
			assert.Equal(t, errors.Code(tt.want.err), errors.Code(err))
		})
	}
}

func TestService_Create(t *testing.T) {
	type args struct {
		userID   string
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
		mock func(dep, args, want)
	}{
		{
			name: "create game successfully",
			args: args{userID: "111", settings: domain.GameSettings{Rows: 6, Columns: 6, BombsNumber: 10}},
			want: want{result: MockGame("111", "xyz", domain.GameStateNew)},
			mock: func(dep dep, args args, want want) {
				dep.rnd.EXPECT().GenerateID().Return(want.result.ID)
				dep.repository.EXPECT().Save(want.result).Return(nil)
			},
		},
		{
			name: "fail at save in repository",
			args: args{userID: "111", settings: domain.GameSettings{Rows: 6, Columns: 6, BombsNumber: 10}},
			want: want{result: domain.Game{}, err: errors.New(apperrors.Internal, apperrors.Internal, "an internal error has occurred", "failed at saving game into repository")},
			mock: func(dep dep, args args, want want) {
				dep.rnd.EXPECT().GenerateID().Return("xyz")
				dep.repository.EXPECT().Save(MockGame(args.userID, "xyz", "")).Return(apperrors.Internal)
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			dep := newDep(t)
			service := game.NewService(dep.rnd, dep.clock, dep.repository)
			tt.mock(dep, tt.args, tt.want)
			result, err := service.Create(tt.args.userID, tt.args.settings)

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

	type args struct {
		userID string
		gameID string
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
		mock func(dep, args, want)
	}{
		{
			name: "game not found",
			args: args{userID: "111", gameID: "xyz", row: 1, column: 1},
			want: want{result: domain.Game{}, err: errors.New(apperrors.NotFound, nil, "game has not been found", "")},
			mock: func(dep dep, args args, want want) {
				dep.repository.EXPECT().Get(args.userID, args.gameID).Return(nil, nil)
			},
		},
		{
			name: "game has already been finished - lost",
			args: args{userID: "111", gameID: "xyz", row: 1, column: 1},
			want: want{result: domain.Game{}, err: errors.New(apperrors.InvalidInput, nil, "game has already finished", "")},
			mock: func(dep dep, args args, want want) {
				game := MockGame("111", "xyz", domain.GameStateLost)
				dep.repository.EXPECT().Get(args.userID, args.gameID).Return(&game, nil)
			},
		},
		{
			name: "game has already been finished - won",
			args: args{userID: "111", gameID: "xyz", row: 1, column: 1},
			want: want{result: domain.Game{}, err: errors.New(apperrors.InvalidInput, nil, "game has already finished", "")},
			mock: func(dep dep, args args, want want) {
				game := MockGame("111", "xyz", domain.GameStateWon)
				dep.repository.EXPECT().Get(args.userID, args.gameID).Return(&game, nil)
			},
		},
		{
			name: "invalid position",
			args: args{userID: "111", gameID: "xyz", row: -100, column: 1},
			want: want{result: domain.Game{}, err: errors.New(apperrors.InvalidInput, nil, "invalid row and column parameters", "")},
			mock: func(dep dep, args args, want want) {
				game := MockGame("111", "xyz", domain.GameStateNew)
				dep.repository.EXPECT().Get(args.userID, args.gameID).Return(&game, nil)
			},
		},
		{
			name: "cell is marked",
			args: args{gameID: "xyz", row: 1, column: 1},
			want: want{result: MockGameWithCell("111", "xyz", domain.GameStateOnGoing, 1, 1, domain.EmptyCellCoveredAndMarked)},
			mock: func(dep dep, args args, want want) {
				dep.repository.EXPECT().Get(args.userID, args.gameID).Return(&want.result, nil)
			},
		},
		{
			name: "reveal first cell successfully",
			args: args{userID: "111", gameID: "xyz", row: 2, column: 2},
			want: want{result: MockGameWithBoard("111", "xyz", domain.GameStateOnGoing, 5, domain.Board{
				{b, e, b, e, e, e},
				{e, b, e, e, e, e},
				{e, e, E, e, b, e},
				{e, e, e, e, e, b},
			}, mockedTime, time.Time{})},
			mock: func(dep dep, args args, want want) {
				game := MockGameWithBoard("111", "xyz", domain.GameStateNew, 5, domain.Board{
					{e, e, e, e, e, e},
					{e, e, e, e, e, e},
					{e, e, e, e, e, e},
					{e, e, e, e, e, e},
				}, mockedTime, time.Time{})
				dep.clock.EXPECT().Now().Return(mockedTime)
				dep.rnd.EXPECT().GenerateN(game.Settings.Rows * game.Settings.Columns).Return([]int{2, 7, 16, 14, 23, 0})
				dep.repository.EXPECT().Get(args.userID, args.gameID).Return(&game, nil)
				dep.repository.EXPECT().Save(want.result).Return(nil)
			},
		},
		{
			name: "reveal cell in cascade successfully",
			args: args{userID: "111", gameID: "xyz", row: 1, column: 1},
			want: want{result: MockGameWithBoard("111", "xyz", domain.GameStateOnGoing, 2, domain.Board{
				{E, E, E, E, E, E},
				{E, E, E, E, E, E},
				{E, E, E, b, E, E},
				{E, E, E, e, e, b},
			}, mockedTime, time.Time{})},
			mock: func(dep dep, args args, want want) {
				game := MockGameWithBoard("111", "xyz", domain.GameStateOnGoing, 2, domain.Board{
					{e, e, e, e, e, e},
					{e, e, e, e, e, e},
					{e, e, e, b, e, e},
					{e, e, e, e, e, b},
				}, mockedTime, time.Time{})
				dep.clock.EXPECT().Now().Return(mockedTime)
				dep.repository.EXPECT().Get(args.userID, args.gameID).Return(&game, nil)
				dep.repository.EXPECT().Save(want.result).Return(nil)
			},
		},
		{
			name: "reveal cell with bomb and lost game",
			args: args{userID: "111", gameID: "xyz", row: 2, column: 3},
			want: want{result: MockGameWithBoard("111", "xyz", domain.GameStateLost, 2, domain.Board{
				{e, e, e, e, e, e},
				{e, e, e, e, e, e},
				{e, e, e, B, e, e},
				{e, e, e, e, e, b},
			}, mockedTime, mockedTime)},
			mock: func(dep dep, args args, want want) {
				game := MockGameWithBoard("111", "xyz", domain.GameStateOnGoing, 2, domain.Board{
					{e, e, e, e, e, e},
					{e, e, e, e, e, e},
					{e, e, e, b, e, e},
					{e, e, e, e, e, b},
				}, mockedTime, time.Time{})
				dep.clock.EXPECT().Now().Return(mockedTime)
				dep.repository.EXPECT().Get(args.userID, args.gameID).Return(&game, nil)
				dep.repository.EXPECT().Save(want.result).Return(nil)
			},
		},
		{
			name: "reveal cell and won game",
			args: args{userID: "111", gameID: "xyz", row: 1, column: 1},
			want: want{result: MockGameWithBoard("111", "xyz", domain.GameStateWon, 2, domain.Board{
				{E, E, E, E, E, E},
				{E, E, E, E, E, E},
				{E, E, E, b, E, E},
				{E, E, E, E, E, b},
			}, mockedTime, mockedTime)},
			mock: func(dep dep, args args, want want) {
				game := MockGameWithBoard("111", "xyz", domain.GameStateOnGoing, 2, domain.Board{
					{e, e, e, e, e, e},
					{e, e, e, e, e, e},
					{e, e, e, b, e, e},
					{e, e, e, E, E, b},
				}, mockedTime, time.Time{})
				dep.clock.EXPECT().Now().Return(mockedTime)
				dep.repository.EXPECT().Get(args.userID, args.gameID).Return(&game, nil)
				dep.repository.EXPECT().Save(want.result).Return(nil)
			},
		},
		{
			name: "fail at save in repository",
			args: args{gameID: "xyz", row: 1, column: 1},
			want: want{result: domain.Game{}, err: errors.New(apperrors.Internal, apperrors.Internal, "an internal error has occurred", "failed at saving game into repository")},
			mock: func(dep dep, args args, want want) {
				game := MockGameWithBoard("111", "xyz", domain.GameStateOnGoing, 2, domain.Board{
					{e, e, e, e, e, e},
					{e, e, e, e, e, e},
					{e, e, e, b, e, e},
					{e, e, e, E, E, b},
				}, mockedTime, time.Time{})
				gameResult := MockGameWithBoard("111", "xyz", domain.GameStateWon, 2, domain.Board{
					{E, E, E, E, E, E},
					{E, E, E, E, E, E},
					{E, E, E, b, E, E},
					{E, E, E, E, E, b},
				}, mockedTime, mockedTime)
				dep.clock.EXPECT().Now().Return(mockedTime)
				dep.repository.EXPECT().Get(args.userID, args.gameID).Return(&game, nil)
				dep.repository.EXPECT().Save(gameResult).Return(apperrors.Internal)
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			dep := newDep(t)
			service := game.NewService(dep.rnd, dep.clock, dep.repository)
			tt.mock(dep, tt.args, tt.want)
			result, err := service.RevealCell(tt.args.userID, tt.args.gameID, tt.args.row, tt.args.column)

			assert.Equal(t, tt.want.result, result)
			if err != nil && tt.want.err != nil {
				assert.Equal(t, tt.want.err.Error(), err.Error())
			}
			assert.Equal(t, errors.Code(tt.want.err), errors.Code(err))
		})
	}
}

// ··· Mocking game primitives ··· //

func MockGame(userID string, gameID string, state string) domain.Game {
	game := domain.Game{
		ID:       gameID,
		UserID:   userID,
		Board:    domain.NewEmptyBoard(6, 6),
		Settings: domain.GameSettings{Rows: 6, Columns: 6, BombsNumber: 10},
		State:    state,
	}

	if state == "" {
		game.State = domain.GameStateNew
	}

	return game
}

func MockGameWithCell(userID, gameID string, state string, row int, column int, cell domain.Cell) domain.Game {
	game := MockGame(userID, gameID, state)
	game.Board.Set(domain.NewPosition(row, column), cell)

	return game
}

func MockGameWithBoard(userID, gameID string, state string, bombsNumber int, board domain.Board, startedAt time.Time, endedAt time.Time) domain.Game {
	return domain.Game{
		ID:       gameID,
		UserID:   userID,
		Board:    board,
		Settings: domain.GameSettings{Rows: len(board), Columns: len(board[0]), BombsNumber: bombsNumber},
		State:    state,
	}
}
