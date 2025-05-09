package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"outbox-go-vinicius/internal/authorization"
	"outbox-go-vinicius/internal/mongodb"
	"strconv"
	"time"
)

const totalMessages = 1

func main() {
	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/?directConnection=true"))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("MongoDB is not reachable:", err)
	}

	log.Println("Database connected.")

	db := client.Database("outbox-db")
	if err := mongodb.InitCollections(ctx, db); err != nil {
		log.Fatal("Failed to init collections:", err)
	}
	repo := authorization.NewRepository(db)
	svc := authorization.NewService(repo)

	startTime := time.Now()

	for i := 0; i < totalMessages; i++ {
		cardID := "card_" + strconv.Itoa(i)
		amount := 1000 + i%1000

		job := authorization.Job{CardID: cardID, Amount: amount}
		if err := svc.AuthorizeCard(ctx, job.CardID, job.Amount); err != nil {
			log.Printf("Failed to authorize job %d: %v", i, err)
		} else {
			log.Printf("Authorized job %d: %+v", i, job)
		}
	}

	elapsedTime := time.Since(startTime)
	log.Printf("All authorizations completed and outbox messages saved! Total messages: %d", totalMessages)
	log.Printf("Time taken: %v", elapsedTime)
}
