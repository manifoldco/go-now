package now

import (
	"fmt"
)

// ClientError represents the error return type
type ClientError interface {
	StatusCode() int
	Code() string
	Message() string
	Error() string
}

type errResponse struct {
	statusCode int
	zeitError  *ZeitError
}

func (e errResponse) StatusCode() int {
	return e.statusCode
}

func (e errResponse) Code() string {
	return e.zeitError.Code
}

func (e errResponse) Message() string {
	return e.zeitError.Message
}

func (e errResponse) Error() string {
	return fmt.Sprintf("%s: %s", e.Code(), e.Message())
}

// NewZeitError construct a new ClientError
func NewZeitError(statusCode int, err *ZeitError) ClientError {
	return errResponse{
		statusCode: statusCode,
		zeitError:  err,
	}
}

type clientError struct {
	message string
}

func (c clientError) StatusCode() int {
	return 0
}

func (c clientError) Code() string {
	return "client_error"
}

func (c clientError) Message() string {
	return c.message
}

func (c clientError) Error() string {
	return fmt.Sprintf("%s: %s", c.Code(), c.Message())
}

// NewError construct a new ClientError
func NewError(err string) ClientError {
	return &clientError{
		message: err,
	}
}
