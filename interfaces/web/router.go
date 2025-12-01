package web

import (
	"github.com/gin-gonic/gin"

	"github.com/qs-lzh/movie-reservation/internal/app"
)

func InitRouter(app *app.App) *gin.Engine {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return router
}
