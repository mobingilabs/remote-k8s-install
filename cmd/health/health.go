package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", health)
	r.Run(":880")
}

func health(c *gin.Context) {
	c.JSON(200, nil)
}
