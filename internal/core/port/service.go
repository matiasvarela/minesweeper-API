package port

import "github.com/matiasvarela/minesweeper-API/internal/core/domain"

type GameService interface {
	Get(id string) (domain.Game, error)
	Create(settings domain.GameSettings) (domain.Game, error)
	MarkCell(id string, row int, column int) (domain.Game, error)
	RevealCell(id string, row int, column int) (domain.Game, error)
}
