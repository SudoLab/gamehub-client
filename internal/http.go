package internal

import (
	"net/http"
	"time"
)

// DoWithRetry performs HTTP request with retry logic
func DoWithRetry(client *http.Client, req *http.Request, maxRetries int) (*http.Response, error) {
	var resp *http.Response
	var err error

	for i := 0; i <= maxRetries; i++ {
		resp, err = client.Do(req)
		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		}

		if i < maxRetries {
			time.Sleep(time.Duration(i+1) * time.Second)
		}
	}

	return resp, err
}
