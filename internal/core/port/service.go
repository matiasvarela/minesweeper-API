package port

import "github.com/matiasvarela/minesweeper-API/internal/core/domain"

type GameService interface {
	Get(userID string, gameID string) (domain.Game, error)
	GetAll(userID string) ([]domain.Game, error)
	Create(userID string, settings domain.GameSettings) (domain.Game, error)
	MarkCell(userID string, gameID string, row int, column int) (domain.Game, error)
	RevealCell(userID string, gameID string, row int, column int) (domain.Game, error)
}
