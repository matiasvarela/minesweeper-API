package domain

const (
	ElementEmpty         = Element('e')
	ElementEmptyMarked   = Element('X')
	ElementEmptyRevealed = Element('E')

	ElementBomb          = Element('b')
	ElementBombMarked    = Element('Y')
	ElementBombRevealed  = Element('B')
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
func (board Board) Set(pos Position, element Element) {
	board[pos.Row][pos.Column] = element
}

// Is returns true if the element in the given position is one of the elements given
func (board Board) Is(pos Position, elements ...Element) bool {
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

			if board.Is(current, ElementBomb) {
				return []Position{}
			}

			if !board.Is(current, ElementEmpty) {
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
			if board.Is(NewPosition(row, column), ElementBomb) {
				board.Set(NewPosition(row, column), ElementEmpty)
			}
		}
	}
}

func (board Board) Count(element Element) int {
	count := 0
	for row := range board {
		for column := range board[0] {
			if board[row][column] == element {
				count ++
			}
		}
	}

	return count
}
