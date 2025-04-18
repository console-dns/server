package errors

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

var HTTP_ERROR = errors.New("HTTP_ERROR")

type HttpStatus interface {
	error
	Code() int
	Message() string
	Writer() http.HandlerFunc
}

type httpError struct {
	error
	statusCode int
	msg        string
	writer     http.HandlerFunc
}

func ErrorRedirect(target string) error {
	return &httpError{
		error:      HTTP_ERROR,
		statusCode: 302,
		msg:        "Redirect",
		writer: func(writer http.ResponseWriter, request *http.Request) {
			http.Redirect(writer, request, target, http.StatusFound)
		},
	}
}

func BadRequestErrorf(msg string, args ...any) error {
	return NewHttpError(http.StatusBadRequest, fmt.Sprintf(msg, args...))
}

func ForbiddenErrorf(msg string, args ...any) error {
	return NewHttpError(http.StatusForbidden, fmt.Sprintf(msg, args...))
}

func NotfoundErrorf(msg string, args ...any) error {
	return NewHttpError(http.StatusNotFound, fmt.Sprintf(msg, args...))
}

func NewHttpError(code int, msg string) error {
	return &httpError{
		error:      HTTP_ERROR,
		statusCode: code,
		msg:        msg,
		writer: func(writer http.ResponseWriter, request *http.Request) {
			http.Error(writer, msg, code)
		},
	}
}

func (h *httpError) Code() int {
	return h.statusCode
}

func (h *httpError) Message() string {
	return h.msg
}

func (h *httpError) Writer() http.HandlerFunc {
	return h.writer
}
