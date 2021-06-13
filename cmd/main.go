package main

import (
	"log"
	"net/http"

	"github.com/tech-pool/takemetotheshop-api/api/registration"
)

const address string = ":8080"

func main() {

	http.HandleFunc("/registration", registration.Handler)

	log.Fatalln(http.ListenAndServe(address, nil))
}
