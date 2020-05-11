package dep

import (
	"github.com/matiasvarela/minesweeper-API/internal/core/port"
	"github.com/matiasvarela/minesweeper-API/internal/handler"
	"github.com/matiasvarela/minesweeper-API/pkg/dynamodbiface"
)

type Dep struct {
	DynamoDB       dynamodbiface.DynamoDB
	GameService    port.GameService
	GameHandler    *handler.GameHandler
	GameRepository port.GameRepository
}
