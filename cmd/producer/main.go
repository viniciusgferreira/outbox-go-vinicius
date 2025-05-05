package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"outbox-go-vinicius/internal/authorization"
	"outbox-go-vinicius/internal/mongodb"
	"strconv"
	"sync"
	"time"
)

const (
	totalMessages = 100_000
	numWorkers    = 1000 // adjust based on your system/MongoDB capacity
	batchSize     = 1000 // Adjust based on memory and performance needs
)

type Job struct {
	cardID string
	amount int
}

func main() {
	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/?directConnection=true"))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB is not reachable:", err)
	}

	log.Println("Database connected.")

	db := client.Database("outbox-db")
	if err := mongodb.InitCollections(ctx, db); err != nil {
		log.Fatal("Failed to init collections:", err)
	}
	repo := authorization.NewRepository(db)
	svc := authorization.NewService(repo)

	jobs := make(chan Job, totalMessages)
	batchJobs := make(chan []authorization.Job, numWorkers) // Channel for batches
	var wg sync.WaitGroup
	var counter int

	// Start worker goroutines
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			var batch []authorization.Job
			for job := range jobs {
				// Add jobs to batch
				batch = append(batch, authorization.Job{CardID: job.cardID, Amount: job.amount})

				// If batch is full, process it
				if len(batch) >= batchSize {
					batchJobs <- batch
					// Reset the batch
					batch = nil
				}
			}

			// Process any remaining jobs in the batch
			if len(batch) > 0 {
				batchJobs <- batch
			}
		}(w)
	}

	// Worker for batch processing and logging the batch insert count
	go func() {
		for batch := range batchJobs {
			if err := svc.AuthorizeBatch(ctx, batch); err != nil {
				log.Printf("Failed to process batch: %v", err)
			} else {
				// Successfully processed the batch
				counter += len(batch)

				// Log the count of the current batch
				log.Printf("Processed a batch of %d messages", len(batch))
			}
		}
	}()

	// Send jobs
	startTime := time.Now()
	for i := 0; i < totalMessages; i++ {
		cardID := "card_" + strconv.Itoa(i)
		amount := 1000 + i%1000
		jobs <- Job{cardID, amount}
	}
	close(jobs)

	// Wait for all workers to finish
	wg.Wait()
	close(batchJobs) // Close batchJobs channel once all jobs are processed

	elapsedTime := time.Since(startTime)

	// Log the results
	log.Printf("All authorizations completed and outbox messages saved! Total messages: %d", counter)
	log.Printf("Time taken: %v", elapsedTime)
}
