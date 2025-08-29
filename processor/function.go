package processor

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"cloud.google.com/go/bigquery"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type PubSubMessage struct {
	Data []byte `json:"data"`
}

type EventPayload struct {
	Message PubSubMessage `json:"message"`
}

type WebhookDataBQ struct {
	EventType string          `bigquery:"event_type"`
	Timestamp string          `bigquery:"event_timestamp"`
	RawData   json.RawMessage `bigquery:"raw_data"`
}

type WebhookPayload struct {
	EventType string          `json:"eventType"`
	Timestamp string          `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}

var (
	projectID = os.Getenv("GCP_PROJECT_ID")
	datasetID = os.Getenv("BIGQUERY_DATASET_ID")
	tableID   = os.Getenv("BIGQUERY_TABLE_ID")
	bqClient  *bigquery.Client
)

func init() {
	functions.CloudEvent("ProcessorFunction", ProcessorFunction)

	if projectID == "" || datasetID == "" || tableID == "" {
		log.Fatalf("Variáveis de ambiente GCP_PROJECT_ID, BIGQUERY_DATASET_ID e BIGQUERY_TABLE_ID devem ser definidas.")
	}

	var err error
	ctx := context.Background()
	bqClient, err = bigquery.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Falha ao criar cliente BigQuery: %v", err)
	}
}

func ProcessorFunction(ctx context.Context, e cloudevents.Event) error {
	log.Printf("Processando CloudEvent...")

	var payload EventPayload
	if err := e.DataAs(&payload); err != nil {
		log.Printf("Erro ao extrair dados do CloudEvent: %v", err)
		return nil
	}

	var webhook WebhookPayload
	if err := json.Unmarshal(payload.Message.Data, &webhook); err != nil {
		log.Printf("Erro ao decodificar payload do webhook: %v", err)
		return nil
	}

	row := WebhookDataBQ{
		EventType: webhook.EventType,
		Timestamp: webhook.Timestamp,
		RawData:   webhook.Data,
	}

	inserter := bqClient.Dataset(datasetID).Table(tableID).Inserter()

	if err := inserter.Put(ctx, &row); err != nil {
		log.Printf("Erro ao inserir dados no BigQuery: %v", err)
		return err
	}

	log.Printf("Dados inseridos no BigQuery com sucesso: EventType=%s", row.EventType)
	return nil
}
