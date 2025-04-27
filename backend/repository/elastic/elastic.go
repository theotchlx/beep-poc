package elastic

import (
	"fmt"

	"beep-poc-backend/dto"
)

type MessageRepository interface {
	Save(message *dto.Message) error
	Get(id string) (*dto.Message, error)
	GetAll() ([]*dto.Message, error)
}

type InMemoryMessageRepository struct {
	messages []*dto.Message
}

func NewInMemoryMessageRepository() *InMemoryMessageRepository {
	return &InMemoryMessageRepository{
		messages: make([]*dto.Message, 0, 10),
	}
}

func (r *InMemoryMessageRepository) GetAll() ([]*dto.Message, error) {
	// Return a copy of the slice
	var messages []*dto.Message
	for _, message := range r.messages {
		messages = append(messages, message)
	}
	return messages, nil
}

func (r *InMemoryMessageRepository) Get(id string) (*dto.Message, error) {
	for _, message := range r.messages {
		if message.ID == id {
			return message, nil
		}
	}
	return nil, fmt.Errorf("Message with ID %s not found", id)
}

func (r *InMemoryMessageRepository) Save(message *dto.Message) error {
	// If message already exists, we replace it (status change / update...)
	for i, e := range r.messages {
		if e.ID == message.ID {
			r.messages[i] = message
			return nil
		}
	}

	// Add in in-memory slice.
	r.messages = append(r.messages, message)
	return nil
}
