package elastic

import (
	"context"
	"encoding/json"
	"fmt"

	"beep-poc-backend/dto"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types/enums/operator"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types/enums/textquerytype"
)

type IMessageRepository interface {
	Save(message *dto.Message) error     // Save a message to the repository (create or update).
	Delete(id string) error              // Delete a message by ID.
	Get(id string) (*dto.Message, error) // Get a message by ID.
	GetPaginated(limit int, offset int) ([]dto.Message, error)
	Search(query string, limit int, offset int) ([]dto.Message, error) // Search for messages based on a query string.
}

const indexName = "messages"

type MessageRepository struct {
	client *elasticsearch.TypedClient
}

func NewMessageRepository(client *elasticsearch.TypedClient) *MessageRepository {
	return &MessageRepository{client: client}
}

func (r *MessageRepository) Save(message *dto.Message) error {
	req := r.client.Index(indexName).
		Request(message).
		Id(message.ID)

	_, err := req.Do(context.Background())
	if err != nil {
		return fmt.Errorf("error indexing document ID=%s: %w", message.ID, err)
	}

	return nil
}

func (r *MessageRepository) Delete(id string) error {
	_, err := r.client.Delete(indexName, id).Do(context.Background())
	if err != nil {
		return fmt.Errorf("error deleting document ID=%s: %w", id, err)
	}
	return nil
}

func (r *MessageRepository) Get(id string) (*dto.Message, error) {
	res, err := r.client.Get(indexName, id).Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error getting document ID=%s: %w", id, err)
	}

	if !res.Found {
		return nil, nil // Document not found
	}

	var message dto.Message
	if err := json.Unmarshal(res.Source_, &message); err != nil {
		return nil, fmt.Errorf("error unmarshalling document source: %w", err)
	}

	return &message, nil
}

func (r *MessageRepository) GetPaginated(limit int, offset int) ([]dto.Message, error) {
	res, err := r.client.Search().
		Index(indexName).
		Request(&search.Request{
			Query: &types.Query{
				MatchAll: &types.MatchAllQuery{},
			},
			From: &offset,
			Size: &limit,
		}).
		Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error executing search query: %w", err)
	}

	messages := make([]dto.Message, len(res.Hits.Hits))
	for i, hit := range res.Hits.Hits {
		if err := json.Unmarshal(hit.Source_, &messages[i]); err != nil {
			return nil, fmt.Errorf("error unmarshalling hit source: %w", err)
		}
	}

	return messages, nil
}

func (r *MessageRepository) Search(query string, limit int, offset int) ([]dto.Message, error) {
	res, err := r.client.Search().Index(indexName).Request(&search.Request{
		Query: &types.Query{
			MultiMatch: &types.MultiMatchQuery{
				Query:    query,
				Fields:   []string{"content"}, // Here we search on one field (content) but could add more.
				Operator: &operator.And,
				Type:     &textquerytype.Phraseprefix, // To match on parts of words (instead of whole words).
			},
		},
		From: &offset,
		Size: &limit,
	}).Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error executing search query: %w", err)
	}

	messages := make([]dto.Message, len(res.Hits.Hits))
	for i, hit := range res.Hits.Hits {
		if err := json.Unmarshal(hit.Source_, &messages[i]); err != nil {
			return nil, fmt.Errorf("error unmarshalling hit source: %w", err)
		}
	}

	return messages, nil
}
