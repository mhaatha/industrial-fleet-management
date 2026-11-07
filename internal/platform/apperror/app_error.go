package apperror

import (
	"errors"
	"net/http"
)

// Code is a stable error code that will be consumed by the client
type Code string

const (
	CodeServiceUnavailable = "SERVICE_UNAVAILABLE"
	CodeInternalError      = "INTERNAL_SERVER_ERROR"
)

type AppError struct {
	Code       Code        `json:"code"`
	Message    string      `json:"message"`
	Details    interface{} `json:"details,omitempty"`
	HTTPStatus int         `json:"-"`
	Err        error       `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// Unwrap so that errors.Is / errors.As can trace the original error
func (e *AppError) Unwrap() error { return e.Err }

// WithErr adds the original error for logging purposes
func (e *AppError) WithErr(err error) *AppError {
	clone := *e
	clone.Err = err
	return &clone
}

// WithDetails adds details to the error
func (e *AppError) WithDetails(d interface{}) *AppError {
	clone := *e
	clone.Details = d
	return &clone
}

// ---------- Constructor ----------
func New(code Code, status int, msg string) *AppError {
	return &AppError{Code: code, HTTPStatus: status, Message: msg}
}

func ServiceUnavailable(msg string) *AppError {
	return New(CodeServiceUnavailable, http.StatusServiceUnavailable, msg)
}

func InternalError(err error) *AppError {
	return &AppError{
		Code:       CodeInternalError,
		Message:    "Something went wrong. Please try again later.",
		HTTPStatus: http.StatusInternalServerError,
		Err:        err,
	}
}

// As is a helper to extract *AppError from error chain
func As(err error) (*AppError, bool) {
	var ae *AppError
	if errors.As(err, &ae) {
		return ae, true
	}
	return nil, false
}
