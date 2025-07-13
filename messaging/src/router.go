package main

import "net/http"

func Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /payments-summary", PaymentsSummaryHandler)
	mux.HandleFunc("POST /purge-payments", PurgeHandler)
	return mux
}
