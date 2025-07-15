package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var defaultProcessorURL = ""
var fallbackProcessorURL = ""

func main() {
	defaultProcessorURL = os.Getenv("DEFAULT_PROCESSOR_URL")
	fallbackProcessorURL = os.Getenv("FALLBACK_PROCESSOR_URL")

	if defaultProcessorURL == "" || fallbackProcessorURL == "" {
		if defaultProcessorURL == "" {
			log.Fatal("DEFAULT_PROCESSOR_URL não está definida")
		}
		if fallbackProcessorURL == "" {
			log.Fatal("FALLBACK_PROCESSOR_URL não está definida")
		}
		os.Exit(1)
	}
	fmt.Println("started")

	RunDB()

	for i := 0; i < 20; i++ {
		go RunWorker()
	}

	go RunMessagingServer()
	go RunHttpServer()

	go RunHealthCkeck()

	select {}
}

func RunHttpServer() {
	addr := flag.String("addr", ":5000", "HTTP network address")

	server := &http.Server{
		Addr:    *addr,
		Handler: Routes(),
	}

	fmt.Println("starting http server", "addr", server.Addr)

	errServer := server.ListenAndServe()
	if errServer != nil {
		os.Exit(1)
	}
}
