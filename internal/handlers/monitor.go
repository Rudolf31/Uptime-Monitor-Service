package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"
	"uptime-monitor/internal/models"
	"uptime-monitor/internal/storage"
)

type MonitorHandler struct {
	storage storage.Storage
	logger  *slog.Logger
}

func NewMonitorHandler(s storage.Storage) *MonitorHandler {
	return &MonitorHandler{
		storage: s,
		logger:  slog.Default(),
	}
}

func (h *MonitorHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req models.CreateMonitorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("invalid create payload", "err", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Interval == 0 || req.Interval < 30 {
		req.Interval = 30
	}
	status := req.Status
	if status == "" {
		status = "pending"
	}

	now := time.Now()
	next := now

	monitor := &models.Monitor{
		URL:          req.URL,
		Interval:     req.Interval,
		Status:       status,
		LastCheck:    &now,
		NextCheck:    &next,
		ResponseTime: nil,
	}

	if err := monitor.Validate(); err != nil {
		slog.Warn("validation failed", "err", err, "url", monitor.URL)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.storage.Create(monitor); err != nil {
		slog.Error("failed to create monitor", "err", err, "url", monitor.URL)
		respondWithError(w, http.StatusInternalServerError, "Failed to create monitor")
		return
	}

	respondWithJSON(w, http.StatusCreated, monitor)
}

func (h *MonitorHandler) List(w http.ResponseWriter, r *http.Request) {
	monitors, _ := h.storage.GetAll()

	resp := make([]models.MonitorResponse, 0, len(monitors))
	for _, m := range monitors {
		resp = append(resp, models.MonitorResponse{
			ID:       m.ID,
			URL:      m.URL,
			Interval: m.Interval,
			Status:   m.Status,
		})
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (h *MonitorHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Invalid monitor ID")
		return
	}

	monitor, err := h.storage.GetByID(id)
	if err != nil {
		if errors.Is(err, storage.ErrMonitorNotFound) {
			respondWithError(w, http.StatusNotFound, "Monitor not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Faild to fetch monitor")
		return
	}

	respondWithJSON(w, http.StatusOK, monitor)
}

func (h *MonitorHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid monitor ID")
		return
	}

	if err := h.storage.Delete(id); err != nil {
		if errors.Is(err, storage.ErrMonitorNotFound) {
			respondWithError(w, http.StatusNotFound, "Monitor not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to delete monitor")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
