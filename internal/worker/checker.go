package worker

import (
	"net/http"
	"time"
)

func CheckURL(url string, client *http.Client) (string, time.Duration) {

	start := time.Now()
	resp, err := client.Get(url)
	duration := time.Since(start)

	status := "down"
	if err == nil {
		defer resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			status = "up"
		}
	}

	return status, duration
}
