package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tech-pool/takemetotheshop-api/api/registration"
)

const address string = ":8080"

func main() {

	router := setupRouter()
	router.Run(address)

	//http.HandleFunc("/registration", registration.Handler)
	// log.Fatalln(http.ListenAndServe(address, nil))
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
