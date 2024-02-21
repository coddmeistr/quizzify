package http

import (
	"errors"
	"net/http"
)

const unknown = "UNKNOWN"

var (
	ErrInvalidJSONBody      = errors.New("invalid json body")
	ErrInternal             = errors.New("internal server error")
	ErrFailedValidation     = errors.New("failed validation")
	ErrMaxLimit             = errors.New("limit exceeded")
	ErrMinLimit             = errors.New("not enough for limit")
	ErrNoRequiredValue      = errors.New("no required value")
	ErrUnknown              = errors.New("error unknown")
	ErrForbidden            = errors.New("forbidden resource")
	ErrInvalidTestType      = errors.New("invalid test type")
	ErrInvalidTestStructure = errors.New("invalid test structure")
	ErrNotFound             = errors.New("not found")
)

var codes map[error]string = map[error]string{
	ErrInvalidJSONBody:      "INVALID_JSON_BODY",
	ErrInternal:             "INTERNAL",
	ErrMaxLimit:             "MAX_LIMIT",
	ErrMinLimit:             "MIN_LIMIT",
	ErrFailedValidation:     "FAILED_VALIDATION",
	ErrNoRequiredValue:      "NO_REQUIRED",
	ErrForbidden:            "FORBIDDEN",
	ErrInvalidTestType:      "INVALID_TEST_TYPE",
	ErrInvalidTestStructure: "INVALID_TEST_STRUCTURE",
	ErrNotFound:             "NOT_FOUND",
	ErrUnknown:              unknown,
}

func ErrorCode(err error) string {
	if code, ok := codes[err]; ok {
		return code
	}
	return unknown
}

func CodeFromError(err error) int {
	switch {
	case errors.Is(err, ErrInvalidJSONBody), errors.Is(err, ErrFailedValidation),
		errors.Is(err, ErrMaxLimit), errors.Is(err, ErrMinLimit),
		errors.Is(err, ErrNoRequiredValue), errors.Is(err, ErrInvalidTestType),
		errors.Is(err, ErrInvalidTestStructure):
		return http.StatusBadRequest
	case errors.Is(err, ErrInternal):
		return http.StatusInternalServerError
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrUnknown):
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
