package receiver

import (
	"fmt"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"io"

	"cloud.google.com/go/pubsub"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

var (
	projectID = os.Getenv("GCP_PROJECT_ID")
	pubsubTopic = os.Getenv("PUBSUB_TOPIC_ID")
	pubsubClient *pubsub.Client
)

func init() {
	functions.HTTP("ReceiverFunction", ReceiverFunction)

	if projectID == "" || pubsubTopic == "" {
		log.Fatalf("Variáveis de ambiente GCP_PROJECT_ID e PUBSUB_TOPIC_ID devem ser definidas.")
	}

	var err error
	ctx := context.Background()
	pubsubClient, err = pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Falha ao criar cliente Pub/Sub: %v", err)
	}
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
// ReceiverFunction recebe o webhook, publica no Pub/Sub e responde rapidamente.
// @Summary Recebe um evento de webhook
// @Description Aceita um payload JSON via HTTP POST, publica no Pub/Sub e retorna uma confirmação.
// @Accept  json
// @Produce  json
// @Param   payload body WebhookPayload true "Payload do Webhook"
// @Success 200 {object} SuccessResponse "Resposta de sucesso"
// @Failure 400 {object} ErrorResponse "Corpo do JSON inválido"
// @Failure 405 {object} ErrorResponse "Método não permitido"
// @Failure 500 {object} ErrorResponse "Erro interno ao publicar mensagem"
// @Router / [post]
func ReceiverFunction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Método não permitido"}`, http.StatusMethodNotAllowed)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Erro ao ler o corpo da requisição: %v", err)
		http.Error(w, `{"error":"Não foi possível ler o corpo da requisição"}`, http.StatusInternalServerError)
		return
	}

	if !json.Valid(bodyBytes) {
		http.Error(w, `{"error":"Corpo do JSON inválido"}`, http.StatusBadRequest)
		return
	}

	log.Printf(">>> Webhook recebido, publicando no tópico %s...", pubsubTopic)

	ctx := context.Background()
	topic := pubsubClient.Topic(pubsubTopic)
	result := topic.Publish(ctx, &pubsub.Message{
		Data: bodyBytes,
	})

	if _, err := result.Get(ctx); err != nil {
		log.Printf("Erro ao publicar mensagem no Pub/Sub: %v", err)
		http.Error(w, fmt.Sprintf(`{"error":"Falha ao publicar mensagem: %v"}`, err), http.StatusInternalServerError)
		return
	}

	log.Printf("Mensagem publicada com sucesso!")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SuccessResponse{Status: "sucesso"})
}