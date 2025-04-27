package api

// This package handles the API methods to the Message service, which itself interfaces with the Message repository.

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"beep-poc-backend/dto"
	"beep-poc-backend/service"
)

// Request body validator.

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

// Message API interface, struct, constructor and methods.

type MessageAPI struct {
	server  *echo.Echo
	service service.MessageService
}

func InitMessageAPI(service service.MessageService) *MessageAPI {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	return &MessageAPI{
		server:  e,
		service: service,
	}
}

func (api *MessageAPI) getAllMessages(c echo.Context) error {
	messages, err := api.service.GetAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if len(messages) == 0 {
		return c.String(http.StatusNotFound, "Message not found")
	}
	return c.JSON(http.StatusOK, messages)
}

func (api *MessageAPI) getMessage(c echo.Context) error {
	// First step is to validate the received request into a DTO.
	getMessage := new(dto.GetMessageRequest)
	if err := c.Bind(getMessage); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(getMessage); err != nil {
		return err
	}

	// Then, we call the service to return its response DTO.
	message, err := api.service.Get(getMessage)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if message == nil {
		return c.String(http.StatusNotFound, "Message not found")
	}

	return c.JSON(http.StatusOK, message)
}

func (api *MessageAPI) createMessage(c echo.Context) error {

	createMessage := new(dto.CreateMessageRequest)
	if err := c.Bind(createMessage); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(createMessage); err != nil {
		return err
	}

	message, err := api.service.Save(createMessage)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, message)
}
