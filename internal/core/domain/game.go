package domain

import "time"

const (
	GameStateNew     = "new"
	GameStateOnGoing = "ongoing"
	GameStateLost    = "lost"
	GameStateWon     = "won"
)

type Game struct {
	ID        string       `json:"id"`
	UserID    string       `json:"user_id"`
	Board     Board        `json:"board"`
	Settings  GameSettings `json:"settings"`
	State     string       `json:"state"`
	StartedAt time.Time    `json:"started_at"`
	EndedAt   time.Time    `json:"ended_at"`
}

type GameSettings struct {
	Rows        int `json:"rows"`
	Columns     int `json:"columns"`
	BombsNumber int `json:"bombs_number"`
}
