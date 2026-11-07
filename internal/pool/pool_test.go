package pool_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/AdPatr3s/pre-trainee-practice/internal/pool"
)

func generateTestFiles(t *testing.T) []pool.Job {
	dir := "testdata"
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	var jobs []pool.Job
	id := 1

	for i := 0; i < 10; i++ {
		path := filepath.Join(dir, fmt.Sprintf("file%d.log", i))
		content := fmt.Sprintf("1.1.1.1 GET /api/test 200 100\n2.2.2.2 GET /home 404 50")
		os.WriteFile(path, []byte(content), 0644)
		jobs = append(jobs, pool.Job{ID: id, Path: path})
		id++
	}

	commentPath := filepath.Join(dir, "comments.log")
	os.WriteFile(commentPath, []byte("# just a comment\n# another comment"), 0644)
	jobs = append(jobs, pool.Job{ID: id, Path: commentPath})
	id++

	jobs = append(jobs, pool.Job{ID: id, Path: filepath.Join(dir, "missing.log")})
	id++

	return jobs
}

func TestRun(t *testing.T) {
	jobs := generateTestFiles(t)
	agg, results := pool.Run(context.Background(), 3, jobs)

	if len(results) != 12 {
		t.Errorf("expected 12 results, got %d", len(results))
	}
	if agg.Files != 12 {
		t.Errorf("Files = %d, want 12", agg.Files)
	}
	if agg.Skipped < 2 { // 2 строки в comments.log
		t.Errorf("Skipped = %d, want at least 2", agg.Skipped)
	}
	if agg.ByStatus[200] != 10 {
		t.Errorf("ByStatus[200] = %d, want 10", agg.ByStatus[200])
	}
	if agg.ByStatus[404] != 10 {
		t.Errorf("ByStatus[404] = %d, want 10", agg.ByStatus[404])
	}
	totalWorkerJobs := 0
	for _, count := range agg.PerWorker {
		totalWorkerJobs += count
	}
	if totalWorkerJobs != 12 {
		t.Errorf("sum of PerWorker jobs = %d, want 12", totalWorkerJobs)
	}
}

func TestRun_NoDeadlock(t *testing.T) {
	jobs := []pool.Job{{ID: 1, Path: "testdata/file0.log"}}

	done := make(chan bool)
	go func() {
		pool.Run(context.Background(), 4, jobs)
		done <- true
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Run deadlocked with more workers than jobs")
	}

}

func TestContextCencel(t *testing.T) {
	jobs := generateTestFiles(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	done := make(chan bool)
	go func(){
		pool.Run(ctx, 3, jobs)
		done <- true
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Run deadlocked on canceled context")
	}
}