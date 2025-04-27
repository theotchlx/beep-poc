package api

// This package handles the API methods to the Message service, which itself interfaces with the Message repository.

import (
	"fmt"
	"net/http"
	"strconv"

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
	service service.IMessageService
}

func InitMessageAPI(service service.IMessageService) *MessageAPI {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	return &MessageAPI{
		server:  e,
		service: service,
	}
}

func (api *MessageAPI) getPaginatedMessages(c echo.Context) error {
	// Parse query parameters
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit < 0 {
		return c.String(http.StatusBadRequest, "Invalid or missing 'limit' query parameter")
	}

	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil || offset < 0 {
		return c.String(http.StatusBadRequest, "Invalid or missing 'offset' query parameter")
	}

	// Create the DTO from the parsed query parameters.
	getMessages := &dto.GetMessagesRequest{
		Limit:  limit,
		Offset: offset,
	}

	// Call the service to return its response DTO.
	messages, err := api.service.GetPaginated(getMessages)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, messages)
}

func (api *MessageAPI) getMessage(c echo.Context) error {
	// First step is to validate and unmarshal the received request into a DTO.
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

func (api *MessageAPI) updateMessage(c echo.Context) error {
	// First step is to validate and unmarshal the received request into a DTO.
	updateMessage := new(dto.UpdateMessageRequest)
	if err := c.Bind(updateMessage); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(updateMessage); err != nil {
		return err
	}

	// Then, we call the service to return its response DTO.
	err := api.service.Update(updateMessage)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (api *MessageAPI) searchMessages(c echo.Context) error {
	// Parse query parameters
	query := c.QueryParam("query")
	if query == "" {
		return c.String(http.StatusBadRequest, "Invalid or missing 'query' query parameter")
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit < 0 {
		return c.String(http.StatusBadRequest, "Invalid or missing 'limit' query parameter")
	}

	// Limit the maximum number of messages to 1000.
	// This is to prevent overloading the server with too many messages at once.
	if limit > 1000 {
		return c.String(http.StatusBadRequest, "limit cannot be greater than 1000")
	}

	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil || offset < 0 {
		return c.String(http.StatusBadRequest, "Invalid or missing 'offset' query parameter")
	}

	// Create the DTO from the parsed query parameters.
	searchMessage := &dto.SearchMessagesRequest{
		Query:  query,
		Limit:  limit,
		Offset: offset,
	}

	fmt.Printf("searchMessage: %+v\n", searchMessage)

	// Call the service to return its response DTO.
	messages, err := api.service.Search(searchMessage)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, messages)
}
