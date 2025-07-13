package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")

	server := &http.Server{
		Addr:    *addr,
		Handler: Routes(),
	}
	fmt.Println("starting server", "addr", server.Addr)

	errServer := server.ListenAndServe()
	if errServer != nil {
		os.Exit(1)
	}
}
