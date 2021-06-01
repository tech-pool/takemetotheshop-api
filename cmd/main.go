package main

import (
	"encoding/json"
	"log"
	"net/http"
)

const port string = ":8080"

type HealthStatus struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

type ServerMessage struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/", handler)
	log.Println("[main] Listening on port ", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func handler(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusNotImplemented)
	writer.Write(structToJson(&ServerMessage{Error: true, Message: "Not implemented yet"}))
}

func healthHandler(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Write(structToJson(&HealthStatus{Status: "UP"}))
}

func structToJson(obj interface{}) []byte {
	bytes, err := json.Marshal(obj)
	if err != nil {
		log.Printf("[structToJson] Unable to marshall object to JSON: %s", err.Error())
	}
	return bytes
}
