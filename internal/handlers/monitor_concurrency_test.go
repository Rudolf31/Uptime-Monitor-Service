package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"uptime-monitor/internal/models"
)

func TestConcurrentCreateMonitors(t *testing.T) {
	server := CreateTestServer()
	t.Cleanup(server.Close)

	const count = 1000

	var wg sync.WaitGroup

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			CreateTestMonitor(t, server.URL)
		}(i)
	}

	wg.Wait()

	resp, _ := http.Get(server.URL + "/monitors")
	var monitors []models.Monitor
	json.NewDecoder(resp.Body).Decode(&monitors)
	if len(monitors) != count {
		t.Fatalf("expected %d monitors, got %d", count, len(monitors))
	}
}

func TestConcurrentGetAndDelete(t *testing.T) {
	server := CreateTestServer()
	t.Cleanup(server.Close)

	const count = 1000

	var wg sync.WaitGroup

	monitor := CreateTestMonitor(t, server.URL)

	for i := 0; i < count/2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			http.Get(server.URL + "/monitor/" + fmt.Sprint(monitor.ID))
		}()
	}

	for i := 0; i < count/2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req, _ := http.NewRequest(
				http.MethodDelete,
				server.URL+"/monitor/"+fmt.Sprint(monitor.ID),
				nil,
			)
			http.DefaultClient.Do(req)
		}()
	}

	wg.Wait()
}
