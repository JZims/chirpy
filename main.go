package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	output := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(output))

}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to zero."))

}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {

	port := "8080"
	mux := http.NewServeMux()
	cfg := apiConfig{}
	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}
	appHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	mux.Handle("/app/", cfg.middlewareMetricsInc(appHandler))

	mux.HandleFunc("GET /healthz", healthzHandler)
	mux.HandleFunc("GET /metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /reset", cfg.handlerReset)

	log.Printf("Server started on port: %v", server.Addr)
	log.Fatal(server.ListenAndServe())

}
