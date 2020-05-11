package main

import (
	"github.com/gin-gonic/gin"
	gameService "github.com/matiasvarela/minesweeper-API/internal/core/service/game"
	"github.com/matiasvarela/minesweeper-API/internal/dep"
	"github.com/matiasvarela/minesweeper-API/internal/handler"
	gameRepo "github.com/matiasvarela/minesweeper-API/internal/repository/game"
	"github.com/matiasvarela/minesweeper-API/pkg/clock"
	"github.com/matiasvarela/minesweeper-API/pkg/random"
	"os"
)

func main() {
	d := initDep()
	router := gin.New()

	routes(router, d)
	run(router)
}

func run(router *gin.Engine) {
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}

func initDep() *dep.Dep {
	d := &dep.Dep{}

	if os.Getenv("ENV") == "production" {
		d.DynamoDB = NewProdDynamoDB()
	} else {
		d.DynamoDB = NewLocalDynamoDB()
	}

	rnd := random.NewRandom()
	clk := clock.New()

	d.GameRepository = gameRepo.NewDynamoDB(d.DynamoDB)
	d.GameService = gameService.NewService(rnd, clk, d.GameRepository)
	d.GameHandler = handler.NewGameHandler(d.GameService)

	return d
}
