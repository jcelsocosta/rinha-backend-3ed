package main

import "net/http"

func Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /payments", PaymentHandler)
	mux.HandleFunc("POST /purge-payments", PurgeHandler)
	mux.HandleFunc("GET /payments-summary", PaymentsSummaryHandler)

	return mux
}
