package throttle_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/AdPatr3s/pre-trainee-practice/internal/throttle"
)

func TestRetries(t *testing.T) {
	var failCount int32 = 1

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("fail") == "1" {
			if atomic.AddInt32(&failCount, -1) >= 0 {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.URL.Query().Get("notfound") == "1" {
			w.WriteHeader(http.StatusNotFound) // 404
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	res := throttle.Do(context.Background(), srv.Client(), srv.URL+"?fail=1", 1, 10, 2)
	if res[0].Attempts != 2 {
		t.Errorf("500 retry: Attempts = %d, want 2", res[0].Attempts)
	}
	if res[0].Status != 200 {
		t.Errorf("500 retry: Status = %d, want 200", res[0].Status)
	}

	res404 := throttle.Do(context.Background(), srv.Client(), srv.URL+"?notfound=1", 1, 10, 2)
	if res404[0].Attempts != 1 {
		t.Errorf("404: Attempts = %d, want 1", res404[0].Attempts)
	}
	if res404[0].Status != 404 {
		t.Errorf("404: Status = %d, want 404", res404[0].Status)
	}
}

func TestRateLimit(t *testing.T) {
	var mu sync.Mutex
	var timestamps []time.Time

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		timestamps = append(timestamps, time.Now())
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := 15
	perSec := 5
	start := time.Now()
	res := throttle.Do(context.Background(), srv.Client(), srv.URL, n, perSec, 0)
	elapsed := time.Since(start)

	if len(res) != n {
		t.Fatalf("len(res) = %d, want %d", len(res), n)
	}

	mu.Lock()
	defer mu.Unlock()
	
	for _, ts := range timestamps {
		count := 0
		for _, otherTs := range timestamps {
			if !otherTs.Before(ts) && otherTs.Before(ts.Add(time.Second)) {
				count++
			}
		}
		if count > perSec {
			t.Errorf("Rate limit exceeded: %d requests in window starting at %v", count, ts)
		}
	}

	if elapsed < 2500*time.Millisecond {
		t.Errorf("Completed too fast: %v. Rate limit likely not working.", elapsed)
	}
}