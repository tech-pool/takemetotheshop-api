package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tech-pool/takemetotheshop-api/api/registration"
)

var port string

func init() {
	port = os.Getenv("PORT")
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
