package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/solrac97gr/cqrs/events"
	"github.com/solrac97gr/cqrs/models"
	"github.com/solrac97gr/cqrs/repository"
	"github.com/solrac97gr/cqrs/search"
)

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/feeds", listFeedsHandler).Methods(http.MethodGet)
	r.HandleFunc("/feeds/search", searchHandler).Methods(http.MethodGet)
	return r
}

func onCreatedFeed(m events.CreatedFeedMessage) {
	feed := models.Feed{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
	}

	if err := search.IndexFeed(context.Background(), feed); err != nil {
		log.Printf("failed to index feed: %v", err)
	}
}

func listFeedsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	feeds, err := repository.ListFeeds(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(feeds); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "missing query parameter", http.StatusBadRequest)
		return
	}

	feeds, err := search.SearchFeed(ctx, query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(feeds); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
