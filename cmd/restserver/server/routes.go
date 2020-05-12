package server

import (
	"github.com/gin-gonic/gin"
	"github.com/matiasvarela/minesweeper-API/internal/dep"
	"net/http"
)

func routes(router *gin.Engine, dependencies *dep.Dep) {
	router.GET("/ping", func(request *gin.Context) {
		request.String(http.StatusOK, "pong")
	})

	router.POST("/users/:user_id/games", dependencies.GameHandler.Create)
	router.GET("/users/:user_id/games", dependencies.GameHandler.GetAll)
	router.GET("/users/:user_id/games/:game_id", dependencies.GameHandler.Get)
	router.PUT("/users/:user_id/games/:game_id/actions/reveal", dependencies.GameHandler.Reveal)
	router.PUT("/users/:user_id/games/:game_id/actions/mark", dependencies.GameHandler.Mark)
}
