package openai

import (
	"fmt"
	"strings"
)

type OpenaiError struct {
	Code      int
	Detail    string
	ErrorType string
}

type GenericError struct {
	Msg string
}

func NewOpenaiError(code int, detail string) *OpenaiError {
	// Choose appropriate error type based on code and detail
	switch code {
	case 401:
		if strings.Contains(detail, "Incorrect API key") {
			return &OpenaiError{Code: code, Detail: detail, ErrorType: "IncorrectApiKey"}
		} else if strings.Contains(detail, "You must be a member of an organization") {
			return &OpenaiError{Code: code, Detail: detail, ErrorType: "NoOrganizationMembership"}
		} else {
			return &OpenaiError{Code: code, Detail: detail, ErrorType: "InvalidAuthentication"}
		}
	case 429:
		if strings.Contains(detail, "exceeded your current quota") {
			return &OpenaiError{Code: code, Detail: detail, ErrorType: "QuotaExceeded"}
		} else {
			return &OpenaiError{Code: code, Detail: detail, ErrorType: "RateLimitExceeded"}
		}
	case 500:
		return &OpenaiError{Code: code, Detail: detail, ErrorType: "ServerError"}
	case 503:
		return &OpenaiError{Code: code, Detail: detail, ErrorType: "EngineOverloaded"}
	default:
		return &OpenaiError{Code: code, Detail: detail, ErrorType: "UnknownError"}
	}
}

func (e *OpenaiError) Error() string {
	// Format error message based on type
	switch e.ErrorType {
	case "InvalidAuthentication":
		return fmt.Sprintf("Error code %d: Invalid Authentication - %s", e.Code, e.Detail)
	case "IncorrectApiKey":
		return fmt.Sprintf("Error code %d: Incorrect API key provided - %s", e.Code, e.Detail)
	case "NoOrganizationMembership":
		return fmt.Sprintf("Error code %d: You must be a member of an organization to use the API - %s", e.Code, e.Detail)
	case "RateLimitExceeded":
		return fmt.Sprintf("Error code %d: Rate limit reached for requests - %s", e.Code, e.Detail)
	case "QuotaExceeded":
		return fmt.Sprintf("Error code %d: You exceeded your current quota, please check your plan and billing details - %s", e.Code, e.Detail)
	case "ServerError":
		return fmt.Sprintf("Error code %d: The server had an error while processing your request - %s", e.Code, e.Detail)
	case "EngineOverloaded":
		return fmt.Sprintf("Error code %d: The engine is currently overloaded, please try again later - %s", e.Code, e.Detail)
	case "UnknownError":
		return fmt.Sprintf("Error code %d: Unknown error - %s", e.Code, e.Detail)
	default:
		return "An unknown error occurred with the OpenAI API."
	}
}

func NewGenericError(msg string) *GenericError {
	return &GenericError{Msg: msg}
}

func (e *GenericError) Error() string {
	return e.Msg
}
