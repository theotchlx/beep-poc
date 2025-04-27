package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"beep-poc-backend/dto"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esapi"
)

type IMessageRepository interface {
	Save(message *dto.Message) error                                    // Save a message to the repository (create or update).
	Get(id string) (*dto.Message, error)                                // Get a message by ID.
	GetAll(id string, limit int, offset int) ([]*dto.Message, error)
	Search(query string, limit int, offset int) ([]*dto.Message, error) // Search for messages based on a query string.
	SearchTotalQuantity(query string) (int, error)                      // Get the total number of messages that match a query string.
}

const indexName = "messages"

type MessageRepository struct {
	client *elasticsearch.Client
}

func NewMessageRepository(client *elasticsearch.Client) *MessageRepository {
	return &MessageRepository{client: client}
}

func (r *MessageRepository) Save(message *dto.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshaling message: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: message.ID,
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), r.client)
	if err != nil {
		return fmt.Errorf("error getting response: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing document ID=%s: %s", message.ID, res.String())
	}

	return nil
}

func (r *MessageRepository) Get(id string) (*dto.Message, error) {
	req := esapi.GetRequest{
		Index:      indexName,
		DocumentID: id,
	}

	res, err := req.Do(context.Background(), r.client)
	if err != nil {
		return nil, fmt.Errorf("error getting response: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		if res.StatusCode == 404 {
			return nil, nil // Document not found
		}
		return nil, fmt.Errorf("error getting document ID=%s: %s", id, res.String())
	}

	var result struct {
		Source dto.Message `json:"_source"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing the response body: %w", err)
	}

	return &result.Source, nil
}

func (r *MessageRepository) GetAll(limit int, offset int) ([]dto.Message, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		"from": offset,
		"size": limit,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("error encoding query: %w", err)
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(context.Background()),
		r.client.Search.WithIndex(indexName),
		r.client.Search.WithBody(&buf),
		r.client.Search.WithTrackTotalHits(true),
		r.client.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting response: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error searching documents: %s", res.String())
	}

	var result struct {
		Hits struct {
			Hits []struct {
				Source dto.Message `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing the response body: %w", err)
	}

	messages := make([]dto.Message, len(result.Hits.Hits))
	for i, hit := range result.Hits.Hits {
		messages[i] = hit.Source
	}

	return messages, nil
}

func (r *MessageRepository) Search(query string, limit int, offset int) ([]*dto.Message, error) {
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"author^2", "content"},
			},
		},
		"from": offset,
		"size": limit,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return nil, fmt.Errorf("error encoding query: %w", err)
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(context.Background()),
		r.client.Search.WithIndex(indexName),
		r.client.Search.WithBody(&buf),
		r.client.Search.WithTrackTotalHits(true),
		r.client.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting response: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error searching documents: %s", res.String())
	}

	var result struct {
		Hits struct {
			Hits []struct {
				Source dto.Message `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing the response body: %w", err)
	}

	messages := make([]*dto.Message, len(result.Hits.Hits))
	for i, hit := range result.Hits.Hits {
		messages[i] = &hit.Source
	}

	return messages, nil
}

func (r *MessageRepository) SearchTotalQuantity(query string) (int, error) {
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"author^2", "content"},
			},
		},
		"size": 0, // We only need the count
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return 0, fmt.Errorf("error encoding query: %w", err)
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(context.Background()),
		r.client.Search.WithIndex(indexName),
		r.client.Search.WithBody(&buf),
		r.client.Search.WithTrackTotalHits(true),
		r.client.Search.WithPretty(),
	)
	if err != nil {
		return 0, fmt.Errorf("error getting response: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return 0, fmt.Errorf("error searching documents: %s", res.String())
	}

	var result struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("error parsing the response body: %w", err)
	}

	return result.Hits.Total.Value, nil
}
