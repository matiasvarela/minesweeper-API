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

// Get retrieves the game belonging to the userID given and with gameID given
func (srv *service) Get(userID string, gameID string) (domain.Game, error) {
	game, err := srv.repository.Get(userID, gameID)
	if err != nil {
		return domain.Game{}, errors.New(apperrors.Internal, err, "an internal error has occurred", "failed at getting game from repository")
	}

	if game == nil {
		return domain.Game{}, errors.New(apperrors.NotFound, nil, "game has not been found", "")
	}

	return *game, nil
}

// Get retrieves all the games belonging to the userID given
func (srv *service) GetAll(userID string) ([]domain.Game, error) {
	games, err := srv.repository.GetAll(userID)
	if err != nil {
		return nil, errors.New(apperrors.Internal, err, "an internal error has occurred", "failed at searching games from repository")
	}

	return games, nil
}

// Create creates a new game for the user and settings given
func (srv *service) Create(userID string, settings domain.GameSettings) (domain.Game, error) {
	game := domain.Game{
		ID:       srv.rnd.GenerateID(),
		UserID:   userID,
		Settings: settings,
		Board:    domain.NewEmptyBoard(settings.Rows, settings.Columns),
		State:    domain.GameStateNew,
	}

	if err := srv.repository.Save(game); err != nil {
		return domain.Game{}, errors.New(apperrors.Internal, err, "an internal error has occurred", "failed at saving game into repository")
	}

	return game, nil
}

// MarkCell mark/unmark the given cell with a flag
func (srv *service) MarkCell(userID string, gameID string, row int, column int) (domain.Game, error) {
	game, err := srv.Get(userID, gameID)
	if err != nil {
		return domain.Game{}, errors.Wrap(err, err.Error())
	}

	pos := domain.NewPosition(row, column)

	if !game.Board.IsValidPosition(pos) {
		return domain.Game{}, errors.New(apperrors.InvalidInput, nil, "invalid row and column parameters", "")
	}

	switch game.Board.Get(pos) {
	case domain.EmptyCellCovered:
		game.Board.Set(pos, domain.EmptyCellCoveredAndMarked)
	case domain.EmptyCellCoveredAndMarked:
		game.Board.Set(pos, domain.EmptyCellCovered)
	case domain.BombCellCovered:
		game.Board.Set(pos, domain.BombCellCoveredAndMarked)
	case domain.BombCellCoveredAndMarked:
		game.Board.Set(pos, domain.BombCellCovered)
	default:
		return game, nil
	}

	if err := srv.repository.Save(game); err != nil {
		return domain.Game{}, errors.New(apperrors.Internal, err, "an internal error has occurred", "failed at saving game into repository")
	}

	return game, nil
}

// RevealCell reveals the given cell and will reveal recursively the adjacent cells if there is no bomb as neighbor
func (srv *service) RevealCell(userID string, gameID string, row int, column int) (domain.Game, error) {
	game, err := srv.Get(userID, gameID)
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
	case domain.EmptyCellCovered:
		if game.State == domain.GameStateNew {
			srv.startGame(&game, pos)
			break
		}

		srv.revealInCascade(&game, pos)

		if game.Board.Count(domain.EmptyCellRevealed) == game.Settings.Rows*game.Settings.Columns-game.Settings.BombsNumber {
			game.State = domain.GameStateWon
			game.EndedAt = srv.clock.Now()
		}
	case domain.BombCellCovered:
		game.Board.Set(pos, domain.BombCellRevealed)
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
	game.Board.Set(pos, domain.EmptyCellRevealed)

	srv.fillBoardWithBombs(game, pos)
}

func (srv *service) fillBoardWithBombs(game *domain.Game, exclude domain.Position) {
	var row, column int
	var bomb domain.Position

	count := 0
	for _, v := range srv.rnd.GenerateN(game.Settings.Rows * game.Settings.Columns) {
		if count == game.Settings.BombsNumber {
			break
		}

		row = v / game.Settings.Columns
		column = v - row*game.Settings.Columns

		bomb = domain.NewPosition(row, column)
		if bomb == exclude {
			continue
		}

		game.Board.Set(bomb, domain.BombCellCovered)
		count++
	}
}

func (srv *service) revealInCascade(game *domain.Game, pos domain.Position) {
	switch game.Board.Get(pos) {
	case domain.EmptyCellCovered:
		game.Board.Set(pos, domain.EmptyCellRevealed)
	}

	for _, neighbor := range game.Board.GetNeighborsIfNoBombs(pos) {
		srv.revealInCascade(game, neighbor)
	}

	return
}
