package authorization

import "time"

type Authorization struct {
	ID        string    `bson:"_id"`
	CardID    string    `bson:"card_id"`
	Amount    int       `bson:"amount"`
	Status    string    `bson:"status"`
	CreatedAt time.Time `bson:"created_at"`
}

type OutboxMessage struct {
	ID             string      `bson:"_id"`
	EventType      string      `bson:"event_type"`
	Payload        interface{} `bson:"payload"`
	TopicName      string      `bson:"topic_name"`
	IdempotencyKey string      `bson:"idempotency_key"`
	CreatedAt      time.Time   `bson:"created_at"`
}
