package main

import (
	"io"
	"net/http"
	"sync"
	"time"
)

type loadResult struct {
	total        int
	statusCounts map[int]int
	failures     int
	duration     time.Duration
}

func runLoadTest(target string, requests, concurrency int) loadResult {
	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	jobs := make(chan struct{})
	results := make(chan int)
	startedAt := time.Now()

	var workers sync.WaitGroup
	for range min(concurrency, requests) {
		workers.Add(1)
		go func() {
			defer workers.Done()
			for range jobs {
				results <- requestStatus(client, target)
			}
		}()
	}

	go func() {
		for range requests {
			jobs <- struct{}{}
		}
		close(jobs)

		workers.Wait()
		close(results)
	}()

	result := loadResult{statusCounts: make(map[int]int)}
	for status := range results {
		result.total++
		if status == 0 {
			result.failures++
			continue
		}
		result.statusCounts[status]++
	}

	result.duration = time.Since(startedAt)
	return result
}

func requestStatus(client *http.Client, target string) int {
	response, err := client.Get(target)
	if err != nil {
		return 0
	}
	defer response.Body.Close()

	_, _ = io.Copy(io.Discard, response.Body)
	return response.StatusCode
}
