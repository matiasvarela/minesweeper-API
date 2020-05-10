package port

import "github.com/matiasvarela/minesweeper-API/internal/core/domain"

type GameRepository interface {
	Get(id string) (*domain.Game, error)
	Save(game domain.Game) error
}