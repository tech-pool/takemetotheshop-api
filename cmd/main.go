package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tech-pool/takemetotheshop-api/api/registration"
	"gopkg.in/yaml.v2"
)

type yamlCfg struct {
	Port string `yaml:"default_port"`
}

var port string

func init() {
	port = os.Getenv("TECHSHOP_PORT")

	if port == "" {
		cfg := yamlCfg{}

		yamlFile, err := ioutil.ReadFile("../configs/.config.yaml")

		if err != nil {
			log.Fatal(err)
		}

		err = yaml.Unmarshal(yamlFile, &cfg)
		if err != nil {
			log.Fatal(err)
		}

		port = cfg.Port

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
