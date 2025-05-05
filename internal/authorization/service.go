package authorization

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Job struct {
	CardID string
	Amount int
}

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo}
}

func (s *Service) AuthorizeCard(ctx context.Context, cardID string, amount int) error {
	auth := Authorization{
		ID:        uuid.NewString(),
		CardID:    cardID,
		Amount:    amount,
		Status:    "authorized",
		CreatedAt: time.Now(),
	}

	outbox := OutboxMessage{
		ID:             uuid.NewString(),
		EventType:      "authorization.created",
		Payload:        auth,
		TopicName:      "authorizations.topic",
		IdempotencyKey: uuid.NewString(),
		CreatedAt:      time.Now(),
	}

	return s.repo.Save(ctx, auth, outbox)
}

func (s *Service) AuthorizeBatch(ctx context.Context, jobs []Job) error {
	auths := make([]Authorization, len(jobs))
	outboxes := make([]OutboxMessage, len(jobs))

	// Generate data for batch insert
	for i, job := range jobs {
		auth := Authorization{
			ID:        uuid.NewString(),
			CardID:    job.CardID,
			Amount:    job.Amount,
			Status:    "authorized",
			CreatedAt: time.Now(),
		}

		outbox := OutboxMessage{
			ID:             uuid.NewString(),
			EventType:      "authorization.created",
			Payload:        auth,
			TopicName:      "authorizations.topic",
			IdempotencyKey: uuid.NewString(),
			CreatedAt:      time.Now(),
		}

		auths[i] = auth
		outboxes[i] = outbox
	}

	// Save in bulk using batch insert
	return s.repo.SaveBatch(ctx, auths, outboxes)
}
