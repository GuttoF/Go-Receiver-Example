package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/GuttoF/Go-Receiver-Example/receiver"
)

func TestReceiverFunction(t *testing.T) {

	t.Run("Deve processar com sucesso uma requisição POST válida", func(t *testing.T) {
		validJSON := `{"eventType":"message.sent","timestamp":"2025-08-29T10:00:00Z","data":{"contactId":"123"}}`
		
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(validJSON))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		receiver.ReceiverFunction(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Handler retornou status code errado: obteve %v, esperava %v", status, http.StatusOK)
		}

		expectedResponse := `{"status":"sucesso"}`
		if strings.TrimSpace(rr.Body.String()) != expectedResponse {
			t.Errorf("Handler retornou corpo inesperado: obteve %v, esperava %v", rr.Body.String(), expectedResponse)
		}
	})

	t.Run("Deve retornar erro 405 para métodos diferentes de POST", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		receiver.ReceiverFunction(rr, req)

		if status := rr.Code; status != http.StatusMethodNotAllowed {
			t.Errorf("Handler retornou status code errado: obteve %v, esperava %v", status, http.StatusMethodNotAllowed)
		}
	})

	t.Run("Deve retornar erro 400 para um JSON malformado", func(t *testing.T) {
		// JSON com uma vírgula extra que o torna inválido
		invalidJSON := `{"eventType":"message.sent",}`
		
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(invalidJSON))
		rr := httptest.NewRecorder()

		receiver.ReceiverFunction(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Handler retornou status code errado: obteve %v, esperava %v", status, http.StatusBadRequest)
		}
	})

	t.Run("Deve retornar erro 400 para uma requisição sem corpo", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rr := httptest.NewRecorder()

		receiver.ReceiverFunction(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Handler retornou status code errado: obteve %v, esperava %v", status, http.StatusBadRequest)
		}
	})
}