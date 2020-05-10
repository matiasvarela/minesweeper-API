package game

import (
	"github.com/matiasvarela/errors"
	"github.com/matiasvarela/minesweeper-API/internal/core/domain"
	"github.com/matiasvarela/minesweeper-API/internal/core/port"
	"github.com/matiasvarela/minesweeper-API/pkg/apperrors"
	"github.com/matiasvarela/minesweeper-API/pkg/random"
)

type service struct {
	rnd        random.Random
	repository port.GameRepository
}

func NewService(rnd random.Random, repository port.GameRepository) *service {
	return &service{rnd: rnd, repository: repository}
}

func (srv *service) Get(id string) (domain.Game, error) {
	game, err := srv.repository.Get(id)
	if err != nil {
		return domain.Game{}, errors.New(apperrors.Internal, err, "an internal error has occurred", "failed at getting game from repository")
	}

	if game == nil {
		return domain.Game{}, errors.New(apperrors.NotFound, nil, "game has not been found", "")
	}

	return *game, nil
}

func (srv *service) Create(settings domain.GameSettings) (domain.Game, error) {
	game := domain.Game{
		ID:       srv.rnd.GenerateID(),
		Settings: settings,
		Board:    domain.NewEmptyBoard(settings.Rows, settings.Columns),
		State:    domain.GameStateNew,
	}

	if err := srv.repository.Save(game); err != nil {
		return domain.Game{}, errors.New(apperrors.Internal, err, "an internal error has occurred", "failed at saving game into repository")
	}

	return game, nil
}

func (srv *service) MarkSquare(id string, row int, column int) (domain.Game, error) {
	game, err := srv.Get(id)
	if err != nil {
		return domain.Game{}, errors.Wrap(err, err.Error())
	}

	pos := domain.NewPosition(row, column)

	switch game.Board.Get(pos) {
	case domain.ElementEmpty:
		game.Board.Set(pos, domain.ElementEmptyMarked)
	case domain.ElementEmptyMarked:
		game.Board.Set(pos, domain.ElementEmpty)
	case domain.ElementBomb:
		game.Board.Set(pos, domain.ElementBombMarked)
	case domain.ElementBombMarked:
		game.Board.Set(pos, domain.ElementBomb)
	default:
		return game, nil
	}

	if err := srv.repository.Save(game); err != nil {
		return domain.Game{}, errors.New(apperrors.Internal, err, "an internal error has occurred", "failed at saving game into repository")
	}

	return game, nil
}

func (srv *service) RevealSquare(id string, row int, column int) (domain.Game, error) {
	return domain.Game{}, nil
}
