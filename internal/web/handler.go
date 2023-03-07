package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/nickzhog/neo_nasa/internal/config"
	"github.com/nickzhog/neo_nasa/internal/service/neo"
	"github.com/nickzhog/neo_nasa/pkg/logging"
)

func (h *handler) getNeoCount(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	dates := r.URL.Query()["dates"]
	parsedDates, err := neo.ParseDates(dates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	count, err := h.storage.CountForDates(ctx, parsedDates)
	if err != nil {
		if errors.Is(err, neo.ErrNoResult) {
			w.WriteHeader(http.StatusNoContent)
			w.Write([]byte("no content"))
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%v", count)
}

type NeoJSON struct {
	NeoCounts []struct {
		Date  string `json:"date"`
		Count int    `json:"count"`
	} `json:"neo_counts"`
}

func (h *handler) upsertNeoCount(w http.ResponseWriter, r *http.Request) {
	var request NeoJSON
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	neos := make([]neo.Neo, len(request.NeoCounts))
	for k, nc := range request.NeoCounts {
		date, err := time.Parse("2006-01-02", nc.Date)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		neos[k] = neo.NewNeo(date, nc.Count)
	}

	err = h.storage.BatchUpdate(r.Context(), neos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "added %v neo", len(request.NeoCounts))
}

type handler struct {
	logger  *logging.Logger
	cfg     *config.Config
	storage neo.Repository
}

func NewHandler(
	logger *logging.Logger,
	cfg *config.Config,
	storage neo.Repository) *handler {
	return &handler{
		logger:  logger,
		cfg:     cfg,
		storage: storage,
	}
}
