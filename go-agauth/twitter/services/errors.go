package services

import (
	"fmt"

	"github.com/go-agauth/twitter/users"
)

// APIError represents a Twitter API Error response
// https://dev.twitter.com/overview/api/response-codes
type APIError struct {
	Errors []ErrorDetail `json:"errors"`
}

// ErrorDetail represents an individual item in an APIError.
type ErrorDetail struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e APIError) Error() string {
	if len(e.Errors) > 0 {
		err := e.Errors[0]
		return fmt.Sprintf("twitter: %d %v", err.Code, err.Message)
	}
	return ""
}

// Empty returns true if empty. Otherwise, at least 1 error message/code is
// present and false is returned.
func (e APIError) Empty() bool {
	if len(e.Errors) == 0 {
		return true
	}
	return false
}

// RelevantError returns any non-nil http-related error (creating the request,
// getting the response, decoding) if any. If the decoded apiError is non-zero
// the apiError is returned. Otherwise, no errors occurred, returns nil.
func RelevantError(httpError error, apiError users.APIError) error {
	if httpError != nil {
		return httpError
	}
	if apiError.Empty() {
		return nil
	}
	return apiError
}
