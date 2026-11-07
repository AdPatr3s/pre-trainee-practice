package fetch_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/AdPatr3s/pre-trainee-practice/internal/fetch"
)

func TestSummarize(t *testing.T) {
	result := []fetch.Result{
		{URL: "a", Status: 200, Duration: 10 * time.Millisecond, Err: nil},
		{URL: "b", Status: 200, Duration: 20 * time.Millisecond, Err: nil},
		{URL: "c", Status: 500, Duration: 5 * time.Millisecond, Err: nil},
		{URL: "d", Status: 0, Duration: 0, Err: errors.New("fail")},
	}

	stats := fetch.Summarize(result)
	if stats.Total != 4 { t.Errorf("Total = %d, want 4", stats.Total) }
	if stats.OK != 2 { t.Errorf("OK = %d, want 2", stats.OK) }
	if stats.Failed != 2 { t.Errorf("Failed = %d, want 2", stats.Failed) }
	if stats.Average != 15 * time.Millisecond { t.Errorf("Average = %v, want 15ms", stats.Average) }
	if stats.FastestURL != "a" { t.Errorf("FastestURL = %s, want a", stats.FastestURL) }
	if stats.SlowestURL != "b" { t.Errorf("SlowestURL = %s, want b", stats.SlowestURL) }
}

func TestFetch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/ok":
			w.WriteHeader(200)
			w.Write([]byte("<html><title>Hello</title></html>"))
		case "/slow":
			time.Sleep(20 * time.Millisecond)
			w.WriteHeader(200)
			w.Write([]byte("<html><title>Slow</title></html>"))
		case "/err":
			w.WriteHeader(500)
		}
	}))
	defer srv.Close()
	urls := []string{srv.URL + "/ok", srv.URL + "/slow", srv.URL + "/err"}
	res := fetch.Fetch(context.Background(), srv.Client(), urls, 2)

	if len(res) != 3 {
		t.Fatalf("len = %d, want 3", len(res))
	}
	if res[0].Status != 200 || res[0].Title != "Hello" || res[0].Err != nil {
		t.Errorf("/ok failed: %+v", res[0])
	}
	if res[2].Status != 500 || res[2].Err != nil {
		t.Errorf("/err failed: %+v", res[2])
	}
}

func TestFetchLimit(t *testing.T) {
	var current int32
	var maxConcurrent int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cur := atomic.AddInt32(&current, 1)
		
		for {
			max := atomic.LoadInt32(&maxConcurrent)
			if cur <= max || atomic.CompareAndSwapInt32(&maxConcurrent, max, cur) {
				break
			}
		}
		
		time.Sleep(30 * time.Millisecond)
		
		atomic.AddInt32(&current, -1)
	}))
	defer srv.Close()

	urls := make([]string, 10)
	for i := range urls {
		urls[i] = srv.URL
	}
	
	fetch.Fetch(context.Background(), srv.Client(), urls, 3)

	if maxConcurrent != 3 {
		t.Errorf("maxConcurrent = %d, want 3", maxConcurrent)
	}
}