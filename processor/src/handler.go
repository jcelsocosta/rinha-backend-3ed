package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PaymentInput struct {
	CorrelationId string  `json:"correlationId"`
	Amount        float64 `json:"amount"`
}

func PaymentHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var paymentInput PaymentInput
	err := json.NewDecoder(r.Body).Decode(&paymentInput)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if paymentInput.CorrelationId == "" || paymentInput.Amount <= 0 {
		http.Error(w, "CorrelationId e Amount são obrigatórios", http.StatusBadRequest)
		return
	}

	err = Publish(paymentInput.CorrelationId, paymentInput.Amount)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func PaymentsSummaryHandler(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	url := fmt.Sprintf("http://messaging:5000/payments-summary?from=%s&to=%s", from, to)

	response, err := http.Get(url)
	if err != nil {
		http.Error(w, "Erro ao fazer requisição para o processador", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
	io.Copy(w, response.Body)
}

func PurgeHandler(w http.ResponseWriter, r *http.Request) {
	url := "http://messaging:5000/purge-payments"

	response, err := http.Post(url, "application/json", nil)
	if err != nil {
		http.Error(w, "Erro ao fazer requisição para o processador", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	io.Copy(w, response.Body)
}
