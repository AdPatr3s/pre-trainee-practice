package throttle

import (
	"context"
	"net/http"
	"sync"
	"time"
)

type Response struct {
	ID        int
	Status    int
	Attempts  int
	Duration  time.Duration
	StartedAt time.Time
	Err       error
}

func Do(ctx context.Context, client *http.Client, url string, n int, perSec int, retries int) []Response {
	if perSec < 1 {
		perSec = 1
	}
	if n < 0 {
		n = 0
	}
	if n == 0 {
		return []Response{}
	}
	if client == nil {
		client = http.DefaultClient
	}

	interval := time.Second / time.Duration(perSec) + time.Millisecond
	ticker := time.NewTicker(interval)

	tokens := make(chan struct{})

	go func() {
		for {
			select {
			case <-ctx.Done():
				return 
			case <-ticker.C:
				tokens <- struct{}{}
			}
		}
	}()
	
	results := make([]Response, n)
	var wg sync.WaitGroup
	wg.Add(n)
	for i:= 0; i < n; i++ {
		go func(id int) {
			defer wg.Done()

			start := time.Now()
			status, attempts, err := doRequest(ctx, client, url, retries, tokens)
			results[id] = Response{
				ID: id,
				Status: status,
				Attempts: attempts,
				Duration: time.Since(start),
				Err: err,
			}

		}(i)
	}
	wg.Wait()
	return results
}

func doRequest(ctx context.Context, client *http.Client, url string, retries int, tokens <-chan struct{}) (status int, attempts int, err error) {
	for attempt := 1; attempt <= retries + 1; attempt++ {
		select {
		case <-ctx.Done():
			return 0, attempt-1, ctx.Err()
		case <-tokens:
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return 0, attempt, err
		}
		resp, err := client.Do(req)
		if err != nil {
			if attempt == retries + 1 {
				return 0, attempt, 	err
			}
			continue
		}
		status = resp.StatusCode
		resp.Body.Close()
		if status >= 500 {
			if attempt == retries + 1 {
				return status, attempt, nil
			}
			continue
		}
		return status, attempt, nil
	}
	return 0, 0, nil
}