package main

import (
	"github.com/matiasvarela/minesweeper-API/cmd/restserver/server"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	server.Start()
}
