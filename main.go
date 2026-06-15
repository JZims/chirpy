package main

import (
	"log"
	"net/http"
)

func main() {
	port := "8080"
	handler := http.NewServeMux()

	server := http.Server{
		Handler: handler,
		Addr:    ":" + port,
	}

	log.Printf("Server started on port: %v", server.Addr)
	log.Fatal(server.ListenAndServe())

}
