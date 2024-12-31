package ilosrv

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BuildServer() *gin.Engine {
	r := gin.Default()

	r.GET("/foo", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]any{
			"message": "hello world!",
		})
	})

	return r
}
