package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	port := flag.String("port", ":8080", "running port")
	flag.Parse()

	router := gin.Default()

	router.GET("/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "hello %s", name)
	})

	router.Run(*port)
}
