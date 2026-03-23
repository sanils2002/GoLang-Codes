package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/bigquery"
)

// checkQueryPendingTime runs a simple, UNIQUE query to bypass the cache
// and returns how long it was pending.
func checkQueryPendingTime(projectID string) (time.Duration, error) {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return 0, fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	// ⭐️ CACHE BUSTING TECHNIQUE ⭐️
	// By adding the current time in nanoseconds as a comment, we ensure
	// the query string is unique for every execution.
	uniqueQuery := fmt.Sprintf("SELECT 1 -- cache_buster_%d", time.Now().UnixNano())
	log.Printf("Running unique query: %s", uniqueQuery)

	q := client.Query(uniqueQuery)
	job, err := q.Run(ctx)
	if err != nil {
		return 0, fmt.Errorf("job.Run: %w", err)
	}

	// The rest of the logic is the same...
	status, err := job.Status(ctx)
	if err != nil {
		return 0, fmt.Errorf("job.Status: %w", err)
	}
	creationTime := status.Statistics.CreationTime

	var startTime time.Time
	for {
		status, err := job.Status(ctx)
		if err != nil {
			return 0, fmt.Errorf("job.Status poll: %w", err)
		}
		if !status.Statistics.StartTime.IsZero() {
			startTime = status.Statistics.StartTime
			break
		}
		time.Sleep(200 * time.Millisecond)
	}

	pendingTime := startTime.Sub(creationTime)
	return pendingTime, nil
}

func main() {
	projectID := "datapipelineproduction" // 👈 Replace with your project ID

	log.Println("Checking BigQuery query pending time with cache-busting...")

	pendingDuration, err := checkQueryPendingTime(projectID)
	if err != nil {
		log.Fatalf("Failed to check query time: %v", err)
	}

	log.Printf("Query was pending for: %s", pendingDuration)

	if pendingDuration > 3*time.Second {
		log.Println("⚠️ High pending time detected! This may indicate slot contention.")
	} else {
		log.Println("✅ Pending time is normal.")
	}
}
