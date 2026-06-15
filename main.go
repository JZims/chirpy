package main

import (
	"log"
	"net/http"
)

func main() {

	handler := http.NewServeMux()

	server := http.Server{
		Handler: handler,
		Addr:    ":8080",
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Error running server: ", err)
	}

	log.Printf("Server started on port: %v", server.Addr)

}
