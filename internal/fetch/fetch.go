package fetch

import (
	"context"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Result struct {
	URL      string
	Status   int    // 0 if no response
	Title    string // first <title>...</title>, trimmed, optional empty
	Duration time.Duration
	Err      error
}

type Stats struct {
	Total      int
	OK         int
	Failed     int
	Average    time.Duration // of successful requests only; 0 if none
	FastestURL string
	SlowestURL string
}

var titleRe = regexp.MustCompile(`(?is)<title>(.*?)</title>`)

func Fetch(ctx context.Context, client *http.Client, urls []string, limit int) []Result {
	if limit < 1 {
		limit = 1
	}
	if client == nil {
		client = http.DefaultClient
	}
	out := make([]Result, len(urls))
	var wg sync.WaitGroup
	sem := make(chan struct{}, limit)
	for i, url := range urls {
		wg.Add(1)
		go func (i int, url string)  {
			defer wg.Done()
			sem <- struct{}{}
			defer func() {<-sem} ()

			start := time.Now()
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if err != nil {
				out[i] = Result{URL: url, Err: err, Duration: time.Since(start)}
				return 
			}
			resp, err := client.Do(req)
			if err != nil {
				out[i] = Result{URL: url, Err: err, Duration: time.Since(start)}
				return
			}
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)
			
			title := ""
			if matches := titleRe.FindSubmatch(body); len(matches) > 1 {
				title = strings.TrimSpace(string(matches[1]))
			}
			out[i] = Result{
				URL:      url,
    			Status:   resp.StatusCode,
    			Title:    title,
    			Duration: time.Since(start),
    			Err:      nil,
			}
		}(i, url)
	}
	wg.Wait()

	return out
}

func Summarize(results []Result) Stats {
	var s Stats
	s.Total = len(results)

	if s.Total == 0 {
		return s
	}

	var totalDuration time.Duration
	var minDuration time.Duration
	var maxDuration time.Duration

	for _, r := range results {
		if r.Err == nil && r.Status >= 200 && r.Status < 300 {
			s.OK++
			totalDuration += r.Duration
			if s.OK == 1 {
				minDuration = r.Duration
				maxDuration = r.Duration
				s.FastestURL = r.URL
				s.SlowestURL = r.URL
			} else {
				if r.Duration < minDuration {
					minDuration = r.Duration
					s.FastestURL = r.URL
				} 
				if r.Duration > maxDuration {
					maxDuration = r.Duration
					s.SlowestURL = r.URL
				}
			}
		} else {
			s.Failed++
		}
	}

	if s.OK > 0 {
		s.Average = totalDuration / time.Duration(s.OK)
	}

	return s
}