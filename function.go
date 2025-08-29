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

func ReceiverFunction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	log.Println(">>> Webhook recebido com sucesso!")

	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Corpo do JSON inválido", http.StatusBadRequest)
		return
	}

	log.Printf("Dados recebidos: %+v", data)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "sucesso"})
}

