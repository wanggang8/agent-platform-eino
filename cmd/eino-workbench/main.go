package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"time"

	"agent-platform-eino/internal/einoapp/bootstrap"
	"agent-platform-eino/internal/einoapp/capabilities"
	"agent-platform-eino/internal/einoapp/execution"
	"agent-platform-eino/internal/einoapp/facts"
	"agent-platform-eino/internal/einoapp/httpapi"
	"agent-platform-eino/internal/einoapp/product"
	"agent-platform-eino/internal/einoapp/store/sqlite"
)

type systemClock struct{}

func (systemClock) Now() time.Time { return time.Now().UTC() }

func main() {
	configPath := flag.String("config", "configs/eino-workbench.local.yaml", "path to M1 workbench config")
	flag.Parse()
	config, err := bootstrap.LoadM1Config(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	repository, err := sqlite.OpenQuery(ctx, config.Database.DSN, sqlite.QueryOptions{BusyTimeout: config.Database.BusyTimeout, PoolSize: config.Database.PoolSize})
	if err != nil {
		log.Fatal(err)
	}
	defer repository.Close()
	source, err := capabilities.LoadNewVulnerabilityFixture(config.Fixture.Path, capabilities.FixtureCase(config.Fixture.Case))
	if err != nil {
		log.Fatal(err)
	}
	registry := capabilities.NewRegistry()
	if err := registry.Register(capabilities.NewVulnerabilityReadCapability(config.Capability.Timeout)); err != nil {
		log.Fatal(err)
	}
	freshness, err := facts.NewFreshnessPolicy(config.Freshness.PolicyVersion, config.Freshness.TTL)
	if err != nil {
		log.Fatal(err)
	}
	service, err := execution.NewM1QueryService(execution.M1QueryServiceConfig{
		Repository: repository, Registry: registry, Source: source, Clock: systemClock{}, Freshness: freshness,
	})
	if err != nil {
		log.Fatal(err)
	}
	projection := product.NewM1Projection(repository)
	handler := httpapi.NewM1Router(httpapi.M1Dependencies{
		Projection: projection,
		Executor:   service,
		Readiness: func(requestContext context.Context) (bool, string) {
			report := repository.Readiness(requestContext)
			return report.Ready, report.Reason
		},
	})
	server := &http.Server{Addr: config.Server.Addr, ReadTimeout: config.Server.ReadTimeout, WriteTimeout: config.Server.WriteTimeout, Handler: handler}
	log.Printf("eino-workbench M1 listening on http://%s", config.Server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
