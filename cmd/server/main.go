package main

import(
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.New()

	router.GET("/ping", func(request *gin.Context){
		request.String(http.StatusOK, "pong")
	})

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}