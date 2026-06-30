package main

import (
	"log"
	"net/http"

	"agent-platform-eino/internal/einoapp/bootstrap"
	"agent-platform-eino/internal/einoapp/httpapi"
)

func main() {
	cfg := bootstrap.LoadConfig()
	server := &http.Server{
		Addr:    cfg.Addr,
		Handler: httpapi.NewRouter(),
	}

	log.Printf("eino-workbench listening on http://%s", cfg.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
