package main

import (
	"flag"
	"log"
	"net/http"

	"agent-platform-eino/internal/einoapp/bootstrap"
	"agent-platform-eino/internal/einoapp/httpapi"
)

func main() {
	configPath := flag.String("config", "configs/eino-workbench.local.yaml", "path to eino-workbench config file")
	flag.Parse()

	cfg, err := bootstrap.LoadConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	server := &http.Server{
		Addr:         cfg.Server.Addr,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		Handler:      httpapi.NewRouter(httpapi.DefaultDependencies()),
	}

	log.Printf("eino-workbench listening on http://%s", cfg.Server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
