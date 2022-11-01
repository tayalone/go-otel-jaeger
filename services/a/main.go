package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	p := "3000"

	pEnv := os.Getenv("PORT")

	if pEnv != "" {
		p = pEnv
	}

	r.Run(":" + p) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
