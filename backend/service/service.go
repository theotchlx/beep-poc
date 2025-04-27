package service

// This package implements service logic to interface with the repositories.

import (
	"time"

	"github.com/google/uuid"

	"beep-poc-backend/dto"
	"beep-poc-backend/repository/elastic"
)

// Message service interface, struct, constructor and methods.

type IMessageService interface {
	Save(request *dto.CreateMessageRequest) (*dto.CreateMessageResponse, error)
	Get(request *dto.GetMessageRequest) (*dto.GetMessageResponse, error)
	GetPaginated(request *dto.GetMessagesRequest) ([]*dto.GetMessageResponse, error)
	Search(request *dto.SearchMessagesRequest) ([]*dto.GetMessageResponse, error)
}

type MessageService struct {
	messageRepository elastic.IMessageRepository
}

func InitMessageService(messageRepository elastic.IMessageRepository) *MessageService {
	return &MessageService{
		messageRepository: messageRepository,
	}
}

func (svc *MessageService) GetPaginated(request *dto.GetMessagesRequest) ([]*dto.GetMessageResponse, error) {
	messages, err := svc.messageRepository.GetPaginated(request.Limit, request.Offset) // Get paginated messages
	if err != nil {
		return nil, err
	}

	var response []*dto.GetMessageResponse
	for _, message := range messages {
		response = append(response, &dto.GetMessageResponse{
			ID:        message.ID,
			Author:    message.Author,
			CreatedAt: message.CreatedAt,
			Content:   message.Content,
		})
	}

	return response, nil
}

func (svc *MessageService) Get(request *dto.GetMessageRequest) (*dto.GetMessageResponse, error) {
	message, err := svc.messageRepository.Get(request.ID)
	if err != nil {
		return nil, err
	}

	if message == nil {
		return nil, nil
	}

	// Return the message object as a DTO
	return &dto.GetMessageResponse{
		ID:        message.ID,
		Author:    message.Author,
		CreatedAt: message.CreatedAt,
		Content:   message.Content,
	}, nil
}

// Enroll a student in a course
func (svc *MessageService) Save(request *dto.CreateMessageRequest) (*dto.CreateMessageResponse, error) {
	/*  1. Save the message in the message repository.
	 *  2. Return the message to the caller.
	 */

	// 1. Save the message in the message repository.
	id := uuid.New().String()
	err := svc.messageRepository.Save(&dto.Message{
		ID:        id,
		Author:    request.Author,
		CreatedAt: time.Now(),
		Content:   request.Content,
	})
	if err != nil {
		return nil, err
	}

	// 2. Return the message to the caller
	return &dto.CreateMessageResponse{
		MessageID: id,
	}, nil
}

func (svc *MessageService) Search(request *dto.SearchMessagesRequest) ([]*dto.GetMessageResponse, error) {
	/*  1. Search for messages in the message repository.
	 *  3. Return the messages and total number of messages to the caller.
	 */
	
	messages, err := svc.messageRepository.Search(request.Query, request.Limit, request.Offset) // Get paginated messages
	if err != nil {
		return nil, err
	}

	var response []*dto.GetMessageResponse
	for _, message := range messages {
		response = append(response, &dto.GetMessageResponse{
			ID:        message.ID,
			Author:    message.Author,
			CreatedAt: message.CreatedAt,
			Content:   message.Content,
		})
	}

	return response, nil
}
