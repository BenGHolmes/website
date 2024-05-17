package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.StaticFile("/", "index.html")

	// Listen and serve on 0.0.0.0:80
	router.Run(":80")
}
