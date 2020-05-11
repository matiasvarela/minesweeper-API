package domain_test

import (
	"github.com/matiasvarela/minesweeper-API/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	e = domain.EmptyCellCovered
	b = domain.BombCellCovered
	E = domain.EmptyCellRevealed
)

func TestNewEmptyBoard(t *testing.T) {
	// Setup
	board := domain.Board{
		{e, e, e, e, e, e},
		{e, e, e, e, e, e},
		{e, e, e, e, e, e},
	}

	// Execute
	result := domain.NewEmptyBoard(len(board), len(board[0]))

	// Verify
	assert.Equal(t, board, result)
}

func TestBoard_Is(t *testing.T) {
	// Setup
	board := domain.Board{
		{e, e, e, e, e, e},
		{e, e, b, e, e, e},
		{e, e, e, e, e, e},
	}

	// Execute
	result := board.Is(domain.NewPosition(1, 2), b)

	// Verify
	assert.True(t, result)
}

func TestBoard_Get(t *testing.T) {
	// Setup
	board := domain.Board{
		{e, e, e, e, e, e},
		{e, e, b, e, e, e},
		{e, e, e, e, e, e},
	}

	// Execute
	result := board.Get(domain.NewPosition(1, 2))

	// Verify
	assert.Equal(t, b, result)
}

func TestBoard_Set(t *testing.T) {
	// Setup
	board := domain.Board{
		{e, e, e, e, e, e},
		{e, e, e, e, e, e},
		{e, e, e, e, e, e},
	}

	// Execute
	board.Set(domain.NewPosition(1, 2), b)

	// Verify
	assert.Equal(t, b, board[1][2])
}

func TestBoard_IsValidPosition(t *testing.T) {
	board := domain.Board{
		{e, e, e, e, b, e},
		{e, b, e, e, e, b},
		{b, e, e, e, E, e},
	}

	type args struct {
		row    int
		column int
	}

	type want struct {
		result bool
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "valid position",
			args: args{row: 1, column: 2},
			want: want{result: true},
		},
		{
			name: "invalid row position",
			args: args{row: len(board), column: 1},
			want: want{result: false},
		},
		{
			name: "invalid column position",
			args: args{row: 1, column: len(board[0])},
			want: want{result: false},
		},
		{
			name: "invalid row position - is negative",
			args: args{row: -1, column: 1},
			want: want{result: false},
		},
		{
			name: "invalid column position - is negative",
			args: args{row: 1, column: -1},
			want: want{result: false},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			got := board.IsValidPosition(domain.NewPosition(tt.args.row, tt.args.column))

			assert.Equal(t, tt.want.result, got)
		})
	}
}

func TestBoard_GetNeighborsIfNoBombs(t *testing.T) {
	board := domain.Board{
		{e, e, e, e, e, e},
		{e, e, e, e, e, e},
		{e, e, e, e, e, e},
	}

	boardWithBombs := domain.Board{
		{e, e, e, e, e, e},
		{b, e, e, e, e, e},
		{e, e, e, e, b, e},
	}

	type args struct {
		board  domain.Board
		row    int
		column int
	}

	type want struct {
		result []domain.Position
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "get neighbors from middle",
			args: args{board: board, row: 1, column: 2},
			want: want{result: []domain.Position{
				{Row: 2, Column: 3}, {Row: 2, Column: 2}, {Row: 2, Column: 1},
				{Row: 1, Column: 3}, {Row: 1, Column: 1},
				{Row: 0, Column: 3}, {Row: 0, Column: 2}, {Row: 0, Column: 1}},
			},
		},
		{
			name: "get neighbors from up-right corner",
			args: args{board: board, row: 0, column: 5},
			want: want{result: []domain.Position{
				{Row: 1, Column: 5}, {Row: 1, Column: 4}, {Row: 0, Column: 4}},
			},
		},
		{
			name: "get neighbors from down-left corner",
			args: args{board: board, row: 2, column: 0},
			want: want{result: []domain.Position{
				{Row: 2, Column: 1}, {Row: 1, Column: 1}, {Row: 1, Column: 0}},
			},
		},
		{
			name: "get empty due to at least one neighbor has a bomb",
			args: args{board: boardWithBombs, row: 1, column: 1},
			want: want{result: []domain.Position{}},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.board.GetNeighborsIfNoBombs(domain.NewPosition(tt.args.row, tt.args.column))

			assert.Equal(t, tt.want.result, got)
		})
	}
}

func TestBoard_HideBombs(t *testing.T) {
	// Setup
	board := domain.Board{
		{e, E, e, b, e, b},
		{e, e, b, e, b, e},
		{b, e, e, b, e, e},
	}

	// Execute
	board.HideBombs()

	// Verify
	assert.Equal(t, domain.Board{
		{e, E, e, e, e, e},
		{e, e, e, e, e, e},
		{e, e, e, e, e, e},
	}, board)
}
