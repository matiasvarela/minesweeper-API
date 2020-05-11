package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/matiasvarela/errors"
	"github.com/matiasvarela/minesweeper-API/internal/core/domain"
	"github.com/matiasvarela/minesweeper-API/internal/core/port"
	"github.com/matiasvarela/minesweeper-API/pkg/apierror"
	"github.com/matiasvarela/minesweeper-API/pkg/apperrors"
	"net/http"
)

type GameHandler struct {
	gameService port.GameService
}

func NewGameHandler(gameService port.GameService) *GameHandler {
	return &GameHandler{gameService: gameService}
}

func (hdl *GameHandler) Get(request *gin.Context) {
	game, err := hdl.gameService.Get(request.Param("id"))
	if err != nil {
		request.AbortWithStatusJSON(apierror.New(err))
		return
	}

	game.Board.HideBombs()

	request.JSON(http.StatusOK, game)
}

func (hdl *GameHandler) Create(request *gin.Context) {
	body := domain.GameSettings{}
	if err := request.BindJSON(&body); err != nil {
		request.AbortWithStatusJSON(apierror.New(errors.New(apperrors.InvalidInput, err, "invalid body", "failed at bind json body")))
		return
	}

	game, err := hdl.gameService.Create(body)
	if err != nil {
		request.AbortWithStatusJSON(apierror.New(err))
		return
	}

	game.Board.HideBombs()

	request.JSON(http.StatusCreated, game)
}

func (hdl *GameHandler) Mark(request *gin.Context) {
	body := struct {
		Row int `json:"row"`
		Column int `json:"column"`
	}{}
	if err := request.BindJSON(&body); err != nil {
		request.AbortWithStatusJSON(apierror.New(errors.New(apperrors.InvalidInput, err, "invalid body", "failed at bind json body")))
		return
	}

	game, err := hdl.gameService.MarkCell(request.Param("id"), body.Row, body.Column)
	if err != nil {
		request.AbortWithStatusJSON(apierror.New(err))
		return
	}

	game.Board.HideBombs()

	request.JSON(http.StatusOK, game)
}

func (hdl *GameHandler) Reveal(request *gin.Context) {
	body := struct {
		Row int `json:"row"`
		Column int `json:"column"`
	}{}
	if err := request.BindJSON(&body); err != nil {
		request.AbortWithStatusJSON(apierror.New(errors.New(apperrors.InvalidInput, err, "invalid body", "failed at bind json body")))
		return
	}

	game, err := hdl.gameService.RevealCell(request.Param("id"), body.Row, body.Column)
	if err != nil {
		request.AbortWithStatusJSON(apierror.New(err))
		return
	}

	game.Board.HideBombs()

	request.JSON(http.StatusOK, game)
}