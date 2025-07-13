package main

import (
	"encoding/json"
	"net/http"
)

func PaymentsSummaryHandler(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	payments := GetPaymentsSummary(from, to)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(payments)
}

func PurgeHandler(w http.ResponseWriter, r *http.Request) {
	err := Purge()

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Failed to purge payments.",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "All payments purged.",
	})
}
