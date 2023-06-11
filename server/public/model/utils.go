package model

import (
	"encoding/base32"
	"encoding/json"

	"github.com/pborman/uuid"
)

type AppError struct {
	Id            string `json:"id"`
	Message       string `json:"message"`               // Message to be display to the end user without debugging information
	DetailedError string `json:"detailed_error"`        // Internal error string to help the developer
	RequestId     string `json:"request_id,omitempty"`  // The RequestId that's also set in the header
	StatusCode    int    `json:"status_code,omitempty"` // The http status code
	Where         string `json:"-"`                     // The function where it happened in the form of Struct.Func
	IsOAuth       bool   `json:"is_oauth,omitempty"`    // Whether the error is OAuth specific
	params        map[string]any
	wrapped       error
}

func NewAppError(where string, id string, params map[string]any, details string, status int) *AppError {
	ap := &AppError{
		Id:            id,
		params:        params,
		Message:       id,
		Where:         where,
		DetailedError: details,
		StatusCode:    status,
		IsOAuth:       false,
	}

	return ap
}

var encoding = base32.NewEncoding("ybndrfg8ejkmcpqxot1uwisza345h769").WithPadding(base32.NoPadding)

func NewId() string {
	return encoding.EncodeToString(uuid.NewRandom())
}

func (err *AppError) ToJSON() string {
	// turn the wrapped error into a detailed message
	detailed := err.DetailedError
	defer func() {
		err.DetailedError = detailed
	}()

	err.wrappedToDetailed()

	b, _ := json.Marshal(err)
	return string(b)
}

func (err *AppError) Unwrap() error {
	return err.wrapped
}

func (er *AppError) Wrap(err error) *AppError {
	er.wrapped = err

	return er
}

func (err *AppError) wrappedToDetailed() {
	if err.wrapped == nil {
		return
	}

	if err.DetailedError != "" {
		err.DetailedError += ", "
	}

	err.DetailedError += err.wrapped.Error()
}
