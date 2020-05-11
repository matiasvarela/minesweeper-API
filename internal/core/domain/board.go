package domain

const (
	EmptyCellCovered          = Cell('e')
	EmptyCellCoveredAndMarked = Cell('X')
	EmptyCellRevealed         = Cell('E')

	BombCellCovered          = Cell('b')
	BombCellCoveredAndMarked = Cell('Y')
	BombCellRevealed         = Cell('B')
)

type Board [][]Cell

type Cell string

type Position struct {
	Row    int
	Column int
}

func NewEmptyBoard(rows int, columns int) Board {
	board := make([][]Cell, rows)
	for i := range board {
		board[i] = make([]Cell, columns)
	}

	for row := range board {
		for column := range board[0] {
			board[row][column] = EmptyCellCovered
		}
	}

	return board
}

func NewPosition(row int, column int) Position {
	return Position{Row: row, Column: column}
}

// Get retrieves the element in the given position
func (board Board) Get(pos Position) Cell {
	return board[pos.Row][pos.Column]
}

// Set put the given element in the given position
func (board Board) Set(pos Position, element Cell) {
	board[pos.Row][pos.Column] = element
}

// Is returns true if the element in the given position is one of the elements given
func (board Board) Is(pos Position, elements ...Cell) bool {
	e := board.Get(pos)
	for _, element := range elements {
		if e == element {
			return true
		}
	}

	return false
}

// IsValidPosition returns true if the given position is within the range of the board; returns false otherwise
func (board Board) IsValidPosition(pos Position) bool {
	return pos.Row >= 0 && pos.Column >= 0 && pos.Row < len(board) && pos.Column < len(board[0])
}

// GetNeighborsIfNoBombs returns the neighbors of the given position or empty if at least one neighbor has a bomb
func (board Board) GetNeighborsIfNoBombs(pos Position) []Position {
	var neighbors []Position
	var current Position

	for di := -1; di <= 1; di++ {
		for dj := -1; dj <= 1; dj++ {
			current = NewPosition(pos.Row-di, pos.Column-dj)

			if current == pos || !board.IsValidPosition(current) {
				continue
			}

			if board.Is(current, BombCellCovered) {
				return []Position{}
			}

			if !board.Is(current, EmptyCellCovered) {
				continue
			}

			neighbors = append(neighbors, current)
		}
	}

	return neighbors
}

// HideBombs replace bombs for empty cells
func (board Board) HideBombs() {
	for row := range board {
		for column := range board[0] {
			if board.Is(NewPosition(row, column), BombCellCovered) {
				board.Set(NewPosition(row, column), EmptyCellCovered)
			}
		}
	}
}

// Count counts the elements of given type
func (board Board) Count(element Cell) int {
	count := 0
	for row := range board {
		for column := range board[0] {
			if board[row][column] == element {
				count++
			}
		}
	}

	return count
}
