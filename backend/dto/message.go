package dto

import (
	"time"
)

type Message struct {
    ID        string    `json:"id"`
    Author    string    `json:"author"`
    CreatedAt time.Time `json:"createdAt"`
    Content   string    `json:"content"`
}

type CreateMessageRequest struct {
    Author    string    `json:"author"`
    Content   string    `json:"content"`
}

type CreateMessageResponse struct {
	MessageID string `json:"messageId"`
}

type GetMessageRequest struct {
	ID string `param:"id" validate:"uuid"`
}

type GetMessageResponse struct {
    ID        string    `json:"id"`
    Author    string    `json:"author"`
    CreatedAt time.Time `json:"createdAt"`
    Content   string    `json:"content"`
}
