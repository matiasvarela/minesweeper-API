package ports

import "github.com/matiasvarela/minesweeper-API/internal/core/domain"

type GameService interface {
	Get(id string) (domain.Game, error)
	Create(settings domain.GameSettings) (domain.Game, error)
	MarkSquare(id string, row int, column int) (domain.Game, error)
	RevealSquare(id string, row int, column int) (domain.Game, error)
}
