package pool

import (
	"context"
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

func Run(ctx context.Context, workers int, jobs []Job) (Aggregate, []FileResult)