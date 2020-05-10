package domain

const (
	ElementEmpty        = Element('E')
	ElementBomb         = Element('B')
	ElementRevealed     = Element('R')
	ElementRevealedBomb = Element('X')
	ElementMark         = Element('Y')
)

type Board [][]Element

type Element rune

type Position struct {
	Row    int
	Column int
}

func NewEmptyBoard(rows int, columns int) Board {
	board := make([][]Element, rows)
	for i := range board {
		board[i] = make([]Element, columns)
	}

	for row := range board {
		for column := range board[0] {
			board[row][column] = ElementEmpty
		}
	}

	return board
}

func NewPosition(row int, column int) Position {
	return Position{Row: row, Column: column}
}

// Get retrieves the element in the given position
func (board Board) Get(pos Position) Element {
	return board[pos.Row][pos.Column]
}

// Set put the given element in the given position
func (board Board) Set(element Element, pos Position) {
	board[pos.Row][pos.Column] = element
}

// Is returns true if the board contains the given element in the given position; returns false otherwise
func (board Board) Is(element Element, pos Position) bool {
	return element == board.Get(pos)
}

// IsValidPosition returns true if the given position is within the range of the board; returns false otherwise
func (board Board) IsValidPosition(pos Position) bool {
	return pos.Row >= 0 && pos.Column >= 0 && pos.Row < len(board) && pos.Column < len(board[0])
}

// Neighbors returns the neighbors of the given position
func (board Board) Neighbors(pos Position) []Position {
	var neighbors []Position
	var current Position

	for di := -1; di <= 1; di++ {
		for dj := -1; dj <= 1; dj++ {
			current = NewPosition(pos.Row-di, pos.Column-dj)

			if current == pos || !board.IsValidPosition(current) {
				continue
			}

			neighbors = append(neighbors, current)
		}
	}

	return neighbors
}

// HideBombs replace bombs for empty squares
func (board Board) HideBombs() {
	for row := range board {
		for column := range board[0] {
			if board.Is(ElementBomb, NewPosition(row, column)) {
				board.Set(ElementEmpty, NewPosition(row, column))
			}
		}
	}
}
