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

	router.GET("/games/:id", dependencies.GameHandler.Get)
	router.POST("/games", dependencies.GameHandler.Create)
	router.PUT("/games/:id/reveal", dependencies.GameHandler.Reveal)
	router.PUT("/games/:id/mark", dependencies.GameHandler.Mark)
}
