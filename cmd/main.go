package  main

import (
	"fmt"
	"github.com/matiasvarela/minesweeper-API/pkg/random"
)

func main() {
	rnd := random.NewRandom()
	rnd.Init()



	println(fmt.Sprintf("%v", rnd.GenerateN(10)))
}