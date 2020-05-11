package dep

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/matiasvarela/minesweeper-API/internal/core/port"
	"github.com/matiasvarela/minesweeper-API/internal/handler"
)

type Dep struct {
	DynamoDB       *dynamodb.DynamoDB
	GameService    port.GameService
	GameHandler    *handler.GameHandler
	GameRepository port.GameRepository
}
