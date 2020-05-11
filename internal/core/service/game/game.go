package game

import (
	"github.com/matiasvarela/errors"
	"github.com/matiasvarela/minesweeper-API/internal/core/domain"
	"github.com/matiasvarela/minesweeper-API/internal/core/port"
	"github.com/matiasvarela/minesweeper-API/pkg/apperrors"
	"github.com/matiasvarela/minesweeper-API/pkg/clock"
	"github.com/matiasvarela/minesweeper-API/pkg/random"
)

type service struct {
	rnd        random.Random
	clock      clock.Clock
	repository port.GameRepository
}

func NewService(rnd random.Random, clock clock.Clock, repository port.GameRepository) *service {
	return &service{rnd: rnd, clock: clock, repository: repository}
}

// Get retrieves the game with id given
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

// Create creates a new game with the settings given
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

// MarkSquare mark/unmark the given square with a flag
func (srv *service) MarkSquare(id string, row int, column int) (domain.Game, error) {
	game, err := srv.Get(id)
	if err != nil {
		return domain.Game{}, errors.Wrap(err, err.Error())
	}

	pos := domain.NewPosition(row, column)

	if !game.Board.IsValidPosition(pos) {
		return domain.Game{}, errors.New(apperrors.InvalidInput, nil, "invalid row and column parameters", "")
	}

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

// RevealSquare reveals the given square and will reveal recursively the adjacent squares if there is no bomb as neighbor
func (srv *service) RevealSquare(id string, row int, column int) (domain.Game, error) {
	game, err := srv.Get(id)
	if err != nil {
		return domain.Game{}, errors.Wrap(err, err.Error())
	}

	if game.State == domain.GameStateLost || game.State == domain.GameStateWon {
		return domain.Game{}, errors.New(apperrors.InvalidInput, nil, "game has already finished", "")
	}

	pos := domain.NewPosition(row, column)

	if !game.Board.IsValidPosition(pos) {
		return domain.Game{}, errors.New(apperrors.InvalidInput, nil, "invalid row and column parameters", "")
	}

	switch game.Board.Get(pos) {
	case domain.ElementEmpty:
		if game.State == domain.GameStateNew {
			srv.startGame(&game, pos)
			break
		}

		srv.revealInCascade(&game, pos)

		if game.Board.Count(domain.ElementEmptyRevealed) == game.Settings.Rows*game.Settings.Columns - game.Settings.BombsNumber{
			game.State = domain.GameStateWon
			game.EndedAt = srv.clock.Now()
		}
	case domain.ElementBomb:
		game.Board.Set(pos, domain.ElementBombRevealed)
		game.State = domain.GameStateLost
		game.EndedAt = srv.clock.Now()
	default:
		return game, nil
	}

	if err := srv.repository.Save(game); err != nil {
		return domain.Game{}, errors.New(apperrors.Internal, err, "an internal error has occurred", "failed at saving game into repository")
	}

	return game, nil
}

func (srv *service) startGame(game *domain.Game, pos domain.Position) {
	game.State = domain.GameStateOnGoing
	game.StartedAt = srv.clock.Now()
	game.Board.Set(pos, domain.ElementEmptyRevealed)

	srv.fillBoardWithBombs(game, pos)
}

func (srv *service) fillBoardWithBombs(game *domain.Game, exclude domain.Position) {
	var row, column int
	var bomb domain.Position

	count := 0
	for _, v := range srv.rnd.GenerateN(game.Settings.BombsNumber + 1) {
		if count == game.Settings.BombsNumber {
			break
		}

		row = v / game.Settings.Columns
		column = v - row*game.Settings.Columns

		bomb = domain.NewPosition(row, column)
		if bomb == exclude {
			continue
		}

		game.Board.Set(bomb, domain.ElementBomb)
		count++
	}
}

func (srv *service) revealInCascade(game *domain.Game, pos domain.Position) {
	switch game.Board.Get(pos) {
	case domain.ElementEmpty:
		game.Board.Set(pos, domain.ElementEmptyRevealed)
	}

	for _, neighbor := range game.Board.GetNeighborsIfNoBombs(pos) {
		srv.revealInCascade(game, neighbor)
	}

	return
}
