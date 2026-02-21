package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"uptime-monitor/internal/models"
	"uptime-monitor/internal/storage"
)

func TestCreateMonitorSuccessful(t *testing.T) {
	server := CreateTestServer()
	t.Cleanup(server.Close)

	body := `{
		"url": "https://music.yandex.ru/",
		"interval": 30
	}`

	resp, err := http.Post(
		server.URL+"/monitors",
		"application/json",
		strings.NewReader(body),
	)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected 201, got %d", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected application/json, got %s", resp.Header.Get("Content-Type"))
	}

	var monitor models.Monitor
	err = json.NewDecoder(resp.Body).Decode(&monitor)
	if err != nil {
		t.Fatal(err)
	}
	if monitor.ID <= 0 {
		t.Errorf("expected positive id, got %d", monitor.ID)
	}
	if monitor.Interval < 30 {
		t.Errorf("expected interval more than 30, got %d", monitor.Interval)
	}
	now := time.Now()
	if monitor.LastCheck.After(now) {
		t.Errorf("expected last check before \"now\", got %v", monitor.LastCheck)
	}
	if monitor.Status != "pending" {
		t.Errorf("expected status \"pending\", got %s", monitor.Status)
	}
}

func TestCreateMonitorFailed(t *testing.T) {
	server := CreateTestServer()
	t.Cleanup(server.Close)

	defer server.Close()

	body := `{
		"url": "//music.yandex.ru/",
		"interval": 0
	}`

	resp, err := http.Post(
		server.URL+"/monitors",
		"application/json",
		strings.NewReader(body),
	)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected application/json, got %s", resp.Header.Get("Content-Type"))
	}
}

func TestGetMonitorSuccessful(t *testing.T) {
	server := CreateTestServer()
	t.Cleanup(server.Close)

	CreateTestMonitor(t, server.URL)
	resp, err := http.Get(server.URL + "/monitors")
	if err != nil {
		t.Fatal(err)
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected application/json, got %s", resp.Header.Get("Content-Type"))
	}

	var monitors []models.Monitor
	err = json.NewDecoder(resp.Body).Decode(&monitors)
	if err != nil {
		t.Fatal(err)
	}
	if len(monitors) == 0 {
		t.Errorf("expected non-empty monitors array")
	}
}

func TestGetMonitorByIdSuccessful(t *testing.T) {
	server := CreateTestServer()
	t.Cleanup(server.Close)

	CreateTestMonitor(t, server.URL)
	resp, err := http.Get(server.URL + "/monitors/1")
	if err != nil {
		t.Fatal(err)
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected application/json, got %s", resp.Header.Get("Content-Type"))
	}

	var monitor models.Monitor
	err = json.NewDecoder(resp.Body).Decode(&monitor)
	if err != nil {
		t.Error(err)
	}

	if monitor.ID <= 0 {
		t.Errorf("expected positive id, got %d", monitor.ID)
	}
	if monitor.Interval < 30 {
		t.Errorf("expected interval more than 30, got %d", monitor.Interval)
	}
	now := time.Now()
	if monitor.LastCheck.After(now) {
		t.Errorf("expected last check before \"now\", got %v", monitor.LastCheck)
	}
}

func TestGetMonitorByIdFailed(t *testing.T) {
	server := CreateTestServer()
	t.Cleanup(server.Close)

	resp, err := http.Get(server.URL + "/monitors/999")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected application/json, got %s", resp.Header.Get("Content-Type"))
	}
}

func TestDeleteMonitorSuccessful(t *testing.T) {
	server := CreateTestServer()
	t.Cleanup(server.Close)

	monitor := CreateTestMonitor(t, server.URL)

	req, err := http.NewRequest(http.MethodDelete, server.URL+"/monitors/"+fmt.Sprint(monitor.ID), nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected 204, got %d", resp.StatusCode)
	}
}

func TestDeleteMonitorFailed(t *testing.T) {
	server := CreateTestServer()
	t.Cleanup(server.Close)

	req, err := http.NewRequest(http.MethodDelete, server.URL+"/monitors/999", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected application/json, got %s", resp.Header.Get("Content-Type"))
	}
}

func CreateTestServer() *httptest.Server {
	storage := storage.NewStorage()
	monitorHandler := NewMonitorHandler(storage)

	mux := http.NewServeMux()

	mux.Handle("POST /monitors", http.HandlerFunc(monitorHandler.Create))
	mux.Handle("GET /monitors", http.HandlerFunc(monitorHandler.List))
	mux.Handle("GET /monitors/{id}", http.HandlerFunc(monitorHandler.Get))
	mux.Handle("DELETE /monitors/{id}", http.HandlerFunc(monitorHandler.Delete))

	server := httptest.NewServer(mux)
	return server
}

func CreateTestMonitor(t *testing.T, serverURL string) models.Monitor {
	body := `{
		"url": "https://httpbin.org/status/200",
		"interval": 30
	}`

	resp, err := http.Post(
		serverURL+"/monitors",
		"application/json",
		strings.NewReader(body),
	)
	if err != nil {
		t.Fatal(err)
	}

	var monitor models.Monitor
	err = json.NewDecoder(resp.Body).Decode(&monitor)
	if err != nil {
		t.Fatal(err)
	}

	return monitor
}
