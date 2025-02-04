package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"
	"github.com/solrac97gr/cqrs/events"
	"github.com/solrac97gr/cqrs/models"
	"github.com/solrac97gr/cqrs/repository"
)

type createFeedRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func newRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/feeds", createFeedHandler).Methods(http.MethodPost)
	return router
}

func createFeedHandler(w http.ResponseWriter, r *http.Request) {
	var req createFeedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdAt := time.Now().UTC()

	id, err := ksuid.NewRandom()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	feed := models.Feed{
		ID:          id.String(),
		Title:       req.Title,
		Description: req.Description,
		CreatedAt:   createdAt,
	}

	if err := repository.InsertFeed(r.Context(), &feed); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := events.PublishCreatedFeed(r.Context(), &feed); err != nil {
		log.Printf("failed to publish created feed: %v", err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feed)
}
