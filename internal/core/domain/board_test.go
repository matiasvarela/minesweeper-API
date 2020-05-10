package domain_test

import (
	"github.com/matiasvarela/minesweeper-API/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	E = domain.ElementEmpty
	B = domain.ElementBomb
	R = domain.ElementRevealed
)

func TestNewEmptyBoard(t *testing.T) {
	// Setup
	board := domain.Board{
		{E, E, E, E, E, E},
		{E, E, E, E, E, E},
		{E, E, E, E, E, E},
	}

	// Execute
	result := domain.NewEmptyBoard(len(board), len(board[0]))

	// Verify
	assert.Equal(t, board, result)
}

func TestBoard_Is(t *testing.T) {
	// Setup
	board := domain.Board{
		{E, E, E, E, E, E},
		{E, E, B, E, E, E},
		{E, E, E, E, E, E},
	}

	// Execute
	result := board.Is(B, domain.NewPosition(1, 2))

	// Verify
	assert.True(t, result)
}

func TestBoard_Get(t *testing.T) {
	// Setup
	board := domain.Board{
		{E, E, E, E, E, E},
		{E, E, B, E, E, E},
		{E, E, E, E, E, E},
	}

	// Execute
	result := board.Get(domain.NewPosition(1, 2))

	// Verify
	assert.Equal(t, B, result)
}

func TestBoard_Set(t *testing.T) {
	// Setup
	board := domain.Board{
		{E, E, E, E, E, E},
		{E, E, E, E, E, E},
		{E, E, E, E, E, E},
	}

	// Execute
	board.Set(B, domain.NewPosition(1, 2))

	// Verify
	assert.Equal(t, B, board[1][2])
}

func TestBoard_IsValidPosition(t *testing.T) {
	board := domain.Board{
		{E, E, E, E, B, E},
		{E, B, E, E, E, B},
		{B, E, E, E, R, E},
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

func TestBoard_Neighbors(t *testing.T) {
	board := domain.Board{
		{E, E, E, E, B, E},
		{E, B, E, E, E, B},
		{B, E, E, E, R, E},
	}

	type args struct {
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
			args: args{row: 1, column: 2},
			want: want{result: []domain.Position{
				{Row:2, Column:3}, {Row:2, Column:2}, {Row:2, Column:1},
				{Row:1, Column:3}, {Row:1, Column:1},
				{Row:0, Column:3}, {Row:0, Column:2}, {Row:0, Column:1}},
			},
		},
		{
			name: "get neighbors from up-right corner",
			args: args{row: 0, column: 5},
			want: want{result: []domain.Position{
				{Row:1, Column:5}, {Row:1, Column:4}, {Row:0, Column:4}},
			},
		},
		{
			name: "get neighbors from down-left corner",
			args: args{row: 2, column: 0},
			want: want{result: []domain.Position{
				{Row:2, Column:1}, {Row:1, Column:1}, {Row:1, Column:0}},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			got := board.Neighbors(domain.NewPosition(tt.args.row, tt.args.column))

			assert.Equal(t, tt.want.result, got)
		})
	}
}

func TestBoard_HideBombs(t *testing.T) {
	// Setup
	board := domain.Board{
		{E, R, E, B, E, B},
		{E, E, B, E, B, E},
		{B, E, E, B, E, E},
	}

	// Execute
	board.HideBombs()

	// Verify
	assert.Equal(t, domain.Board{
		{E, R, E, E, E, E},
		{E, E, E, E, E, E},
		{E, E, E, E, E, E},
	}, board)
}