package port

import "github.com/matiasvarela/minesweeper-API/internal/core/domain"

//go:generate mockgen -source=repository.go -destination=../../../mock/repository.go -package=mock

type GameRepository interface {
	Get(userID string, gameID string) (*domain.Game, error)
	GetAll(userID string) ([]domain.Game, error)
	Save(game domain.Game) error
}
