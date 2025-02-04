package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	"github.com/solrac97gr/cqrs/database"
	"github.com/solrac97gr/cqrs/events"
	"github.com/solrac97gr/cqrs/repository"
	"github.com/solrac97gr/cqrs/search"
)

type Config struct {
	PostgresDB           string `envconfig:"POSTGRES_DB" required:"true"`
	PostgresUser         string `envconfig:"POSTGRES_USER" required:"true"`
	PostgresPassword     string `envconfig:"POSTGRES_PASSWORD" required:"true"`
	NatsAddress          string `envconfig:"NATS_ADDRESS" required:"true"`
	ElasticsearchAddress string `envconfig:"ELASTICSEARCH_ADDRESS" required:"true"`
}

func main() {
	var cfg Config

	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("%v", err)
	}

	addr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		"postgres",
		cfg.PostgresDB,
	)

	repo, err := database.NewPostgresRepository(addr)
	if err != nil {
		log.Fatalf("%v", err)
	}
	repository.SetRepository(repo)
	defer repository.Close()

	es, err := search.NewElasticSearchRepository(fmt.Sprintf("http://%s", cfg.ElasticsearchAddress))
	if err != nil {
		log.Fatalf("%v", err)
	}
	search.SetSearchRepository(es)
	defer search.Close()

	n, err := events.NewNats(fmt.Sprintf("nats://%s", cfg.NatsAddress))
	if err != nil {
		log.Fatalf("%v", err)
	}

	err = n.OnCreateFeed(onCreatedFeed)
	if err != nil {
		log.Fatalf("%v", err)
	}

	events.SetEventStore(n)
	defer events.Close()

	router := newRouter()
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("%v", err)
	}
}
