package routes

import (
	"sorteio-daily/controllers"

	"github.com/gin-gonic/gin"
)

func HandleRequests() {
	r := gin.Default()
	r.POST("/pessoa", controllers.CriaNovaPessoa)
	r.GET("/pessoas", controllers.ListPessoas)
	r.POST("/pessoas/sort", controllers.SorteiaPessoa)
	r.Run(":8000")
}
