package domain

import "time"

const (
	GameStateNew     = "new"
	GameStateOnGoing = "ongoing"
	GameStateLost    = "lost"
	GameStateWon     = "won"
)

type Game struct {
	ID        string       `json:"game"`
	Board     Board        `json:"board"`
	Settings  GameSettings `json:"settings"`
	State     string       `json:"state"`
	StartedAt time.Time    `json:"started_at"`
}

type GameSettings struct {
	Rows        int `json:"rows"`
	Columns     int `json:"columns"`
	BombsNumber int `json:"bombs_number"`
}

func (game *Game) HideBombs() *Game {
	for row := range game.Board {
		for column := range game.Board[0] {
			if game.Board.Is(ElementBomb, NewPosition(row, column)) {
				game.Board.Set(ElementEmpty, NewPosition(row, column))
			}
		}
	}

	return game
}
