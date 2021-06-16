package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tech-pool/takemetotheshop-api/api/registration"
)

// to be put in a config file
const default_address string = ":8080"

var port string

func init() {
	port = os.Getenv("TECHSHOP_PORT")

	if port == "" {
		port = default_address
	}
}

func main() {

	router := setupRouter()
	router.Run(port)

}

func setupRouter() *gin.Engine {

	///Engine  returned by Default() is the framework's instance,
	//it contains the muxer, middleware and configuration settings.
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello Gin")
	})

	r.GET("/registration", registration.Handler)

	return r

}
