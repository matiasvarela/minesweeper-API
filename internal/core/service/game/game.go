package game

import "github.com/matiasvarela/minesweeper-API/internal/core/domain"

type service struct{}

func NewService() *service {
	return &service{}
}

func (srv *service) Get(id string) (domain.Game, error) {
	return domain.Game{}, nil
}

func (srv *service) Create(settings domain.GameSettings) (domain.Game, error) {
	return domain.Game{}, nil
}

func (srv *service) MarkSquare(id string, row int, column int) (domain.Game, error) {
	return domain.Game{}, nil
}

func (srv *service) RevealSquare(id string, row int, column int) (domain.Game, error) {
	return domain.Game{}, nil
}
