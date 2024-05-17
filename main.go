package main

import (
	"os"

	"github.com/gin-gonic/gin"

	// This autoloads the .env file
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	router := gin.Default()
	router.StaticFile("/", "index.html")

	// Read SSL certificate and key file paths from environment variables
	certFile := os.Getenv("SSL_CERT_FILE")
	keyFile := os.Getenv("SSL_KEY_FILE")

	// Listen and serve
	router.RunTLS(":443", certFile, keyFile)
}
