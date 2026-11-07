package pool

import (
	"context"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/AdPatr3s/pre-trainee-practice/internal/accesslog"
)

type Job struct {
	ID   int
	Path string
}

type FileResult struct {
	JobID   int
	Path    string
	Parsed  int
	Skipped int
	Err     error
	Elapsed time.Duration
}

type Aggregate struct {
	Files     int
	Parsed    int
	Skipped   int
	ByStatus  map[int]int
	TopPaths  []accesslog.PathCount // top 5, same tie-break as Stage 2
	PerWorker map[int]int           // worker id -> jobs completed
}

func Run(ctx context.Context, workers int, jobs []Job) (Aggregate, []FileResult) {
	if workers < 1 {
		workers = 1
	}
	jobsCh := make(chan Job)
	var wg sync.WaitGroup
	wg.Add(workers)

	var mu sync.Mutex
	agg := Aggregate{
		ByStatus: make(map[int]int),
		PerWorker: make(map[int]int),
	}
	var result []FileResult
	globalPathFreq := make(map[string]int)

	for i := 1; i<=workers; i++ {
		go func(i int){
			defer wg.Done()
			for job := range jobsCh {
				start := time.Now()
				data, err := os.ReadFile(job.Path)
				res := FileResult{JobID: job.ID, Path: job.Path, Elapsed: time.Since(start)}
				if err != nil {
					res.Err = err
				} else {
					sum := accesslog.Summarize(string(data), 1000)
					res.Parsed = sum.Parsed
					res.Skipped = sum.Skipped

					mu.Lock()

					for status, count := range sum.ByStatus {
						agg.ByStatus[status] += count
					}
					for _, pc := range sum.TopPaths {
						globalPathFreq[pc.Path] += pc.Count
					}
					mu.Unlock()

				}
				mu.Lock()
				result = append(result, res)
				agg.PerWorker[i]++
				mu.Unlock()
			}
			
		}(i)
	}
	go func() {
		for _, i := range jobs {
			select {
			case <- ctx.Done():
				close(jobsCh)
				return 
			case jobsCh <- i:
			}
		}
		close(jobsCh)
	}()

	wg.Wait()

	agg.Files = len(result)
	for _, r := range result {
		agg.Parsed += r.Parsed
		agg.Skipped += r.Skipped
	}

	pairs := make([]accesslog.PathCount, 0, len(globalPathFreq))
	for p, c := range globalPathFreq {
		pairs = append(pairs, accesslog.PathCount{Path: p, Count: c})
	}

	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].Count != pairs[j].Count {
			return pairs[i].Count > pairs[j].Count
		}
		return pairs[i].Path < pairs[j].Path
	})

	if len(pairs) > 5 {
		pairs = pairs[:5]
	}
	agg.TopPaths = pairs

	return agg, result
}