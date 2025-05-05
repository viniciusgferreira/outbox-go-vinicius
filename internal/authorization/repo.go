package authorization

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	authCol   *mongo.Collection
	outboxCol *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		authCol:   db.Collection("authorizations"),
		outboxCol: db.Collection("outbox_messages"),
	}
}

func (r *Repository) Save(ctx context.Context, auth Authorization, outbox OutboxMessage) error {
	session, err := r.authCol.Database().Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		if _, err := r.authCol.InsertOne(sessCtx, auth); err != nil {
			return nil, err
		}
		if _, err := r.outboxCol.InsertOne(sessCtx, outbox); err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}

func (r *Repository) SaveBatch(ctx context.Context, auths []Authorization, outboxes []OutboxMessage) error {
	// Perform a bulk insert for both Authorizations and OutboxMessages
	// Prepare bulk write operations for Authorization and Outbox
	authOps := make([]mongo.WriteModel, len(auths))
	outboxOps := make([]mongo.WriteModel, len(outboxes))

	for i, auth := range auths {
		authOps[i] = mongo.NewInsertOneModel().SetDocument(auth)
	}

	for i, outbox := range outboxes {
		outboxOps[i] = mongo.NewInsertOneModel().SetDocument(outbox)
	}

	// Execute the bulk writes
	bulkWriteOpts := options.BulkWrite().SetOrdered(false) // Set unordered for better performance
	_, err := r.authCol.BulkWrite(ctx, authOps, bulkWriteOpts)
	if err != nil {
		return err
	}

	_, err = r.outboxCol.BulkWrite(ctx, outboxOps, bulkWriteOpts)
	if err != nil {
		return err
	}

	return nil
}
