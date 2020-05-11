package server

import "github.com/gin-gonic/gin"

func Start() {
	dependencies := initDependencies()
	router := gin.New()

	routes(router, dependencies)
	run(router)
}

func run(router *gin.Engine) {
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}