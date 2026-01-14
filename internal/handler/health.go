package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthHandler struct {
    db *pgxpool.Pool
}

func NewHealthHandler(db *pgxpool.Pool) *HealthHandler {
    return &HealthHandler{db: db}
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
    health := struct {
        Status    string `json:"status"`
        Timestamp string `json:"timestamp"`
        Database  string `json:"database"`
    }{
        Timestamp: time.Now().UTC().Format(time.RFC3339),
    }

    // Check database connection
    if err := h.db.Ping(r.Context()); err != nil {
        health.Status = "unhealthy"
        health.Database = "disconnected"
        w.WriteHeader(http.StatusServiceUnavailable)
    } else {
        health.Status = "healthy"
        health.Database = "connected"
        w.WriteHeader(http.StatusOK)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(health)
}
