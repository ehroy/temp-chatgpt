package repository

import (
	"context"
	"sort"
	"strings"
	"sync"

	"emailchatgpt/internal/model"
)

type Repository interface {
	ListMessages(ctx context.Context, email string) ([]model.EmailMessage, error)
}

type MemoryRepository struct {
	mu       sync.RWMutex
	messages []model.EmailMessage
}

func NewMemoryRepository(messages []model.EmailMessage) *MemoryRepository {
	return &MemoryRepository{messages: append([]model.EmailMessage(nil), messages...)}
}

func (r *MemoryRepository) ListMessages(_ context.Context, email string) ([]model.EmailMessage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	needle := strings.ToLower(strings.TrimSpace(email))
	filtered := make([]model.EmailMessage, 0)
	for _, message := range r.messages {
		if needle == "" || strings.ToLower(message.Recipient) == needle {
			filtered = append(filtered, message)
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].ReceivedAt.After(filtered[j].ReceivedAt)
	})
	return filtered, nil
}
