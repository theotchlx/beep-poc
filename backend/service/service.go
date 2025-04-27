package service

// This package implements service logic to interface with the repositories.

import (
	"time"

	"github.com/google/uuid"

	"beep-poc-backend/dto"
	"beep-poc-backend/repository/elastic"
)

// Message service interface, struct, constructor and methods.

type MessageService interface {
	Save(request *dto.CreateMessageRequest) (*dto.CreateMessageResponse, error)
	Get(request *dto.GetMessageRequest) (*dto.GetMessageResponse, error)
	GetAll() ([]*dto.GetMessageResponse, error)
}

type ApiMessageService struct {
	messageRepository elastic.MessageRepository
}

func InitMessageService(messageRepository elastic.MessageRepository) *ApiMessageService {
	return &ApiMessageService{
		messageRepository: messageRepository,
	}
}

func (api *ApiMessageService) GetAll() ([]*dto.GetMessageResponse, error) {
	messages, err := api.messageRepository.GetAll()
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

func (api *ApiMessageService) Get(request *dto.GetMessageRequest) (*dto.GetMessageResponse, error) {
	message, err := api.messageRepository.Get(request.ID)
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
func (api *ApiMessageService) Save(request *dto.CreateMessageRequest) (*dto.CreateMessageResponse, error) {
	/*  1. Save the message in the message repository.
	 *  2. Return the message to the caller.
	 */

	// 1. Save the message in the message repository.
	id := uuid.New().String()
	err := api.messageRepository.Save(&dto.Message{
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
