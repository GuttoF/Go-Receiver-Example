package main

import (
	"log"
	"net/http"
	"receiver/receiver"
)

func main() {
	http.HandleFunc("/", receiver.ReceiverFunction)
	log.Println("Servidor rodando em http://localhost:8080 🚀")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

