package domain

import "time"

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
