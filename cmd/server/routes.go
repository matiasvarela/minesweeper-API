package main

import (
	"github.com/gin-gonic/gin"
	"github.com/matiasvarela/minesweeper-API/internal/dep"
	"net/http"
)

func routes(router *gin.Engine, d *dep.Dep) {
	router.GET("/ping", func(request *gin.Context) {
		request.String(http.StatusOK, "pong")
	})

	router.GET("/games/:id", d.GameHandler.Get)
	router.POST("/games", d.GameHandler.Create)
	router.PUT("/games/:id/reveal", d.GameHandler.Reveal)
	router.PUT("/games/:id/mark", d.GameHandler.Mark)
}
