package main

import (
	"log"
	"net/http"

	_ "github.com/GuttoF/Go-Receiver-Example/docs"
	"github.com/GuttoF/Go-Receiver-Example/receiver"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Receptor de Webhooks API
// @version 1.0
// @description Servidor local para a API de recebimento de webhooks e sua documentação.
// @host localhost:8080
// @BasePath /
func main() {
	http.HandleFunc("/", receiver.ReceiverFunction)

	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	
	log.Println("Servidor rodando em http://localhost:8080 🚀")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

