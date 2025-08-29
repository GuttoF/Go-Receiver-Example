package receiver

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)


func init() {
	functions.HTTP("ReceiverFunction", ReceiverFunction)
}

type WebhookPayload struct {
	EventType string                 `json:"eventType" example:"message.sent"`
	Timestamp string                 `json:"timestamp" example:"2025-08-29T10:00:00Z"`
	Data      map[string]interface{} `json:"data"`
}

type SuccessResponse struct {
	Status string `json:"status" example:"sucesso"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"Método não permitido"`
}


// @title Receptor de Webhooks API
// @version 1.0
// @description API para receber eventos de webhook. Esta função é projetada para ser um gateway robusto e escalável, recebendo notificações, validando-as rapidamente e encaminhando para processamento assíncrono.
// @contact.name Suporte da API
// @contact.url http://www.seusite.com/support
// @contact.email support@seusite.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
// @schemes http
// ReceiverFunction é o handler que recebe e processa os webhooks.
// @Summary Recebe um evento de webhook
// @Description Aceita um payload JSON via HTTP POST, valida-o e retorna uma confirmação de sucesso.
// @Accept  json
// @Produce  json
// @Param   payload body WebhookPayload true "Payload do Webhook"
// @Success 200 {object} SuccessResponse "Resposta de sucesso"
// @Failure 400 {object} ErrorResponse "Corpo do JSON inválido"
// @Failure 405 {object} ErrorResponse "Método não permitido"
// @Router / [post]
func ReceiverFunction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Método não permitido"}`, http.StatusMethodNotAllowed)
		return
	}
	log.Println(">>> Webhook recebido com sucesso!")

	var payload WebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, `{"error":"Corpo do JSON inválido"}`, http.StatusBadRequest)
		return
	}

	log.Printf("Dados recebidos: %+v", payload)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SuccessResponse{Status: "sucesso"})
}

