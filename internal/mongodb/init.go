package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitCollections(ctx context.Context, db *mongo.Database) error {
	err := createCollectionIfNotExists(ctx, db, "authorizations")
	if err != nil {
		return fmt.Errorf("error creating authorizations collection: %w", err)
	}
	err = createCollectionIfNotExists(ctx, db, "outbox_messages")
	if err != nil {
		return fmt.Errorf("error creating outbox_messages collection: %w", err)
	}
	outboxCol := db.Collection("outbox_messages")
	_, err = outboxCol.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "created_at", Value: 1}},
	})
	if err != nil {
		return fmt.Errorf("error creating index: %w", err)
	}
	return nil
}

func createCollectionIfNotExists(ctx context.Context, db *mongo.Database, name string) error {
	collections, err := db.ListCollectionNames(ctx, bson.M{"name": name})
	if err != nil {
		return err
	}
	if len(collections) == 0 {
		return db.CreateCollection(ctx, name)
	}
	return nil
}
