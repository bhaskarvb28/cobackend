package shared

type APIError struct {
	StatusCode int
	Message    string
}

// func (e *APIError) Error() string {
// 	return e.Message
// }

func NewAPIError(
	statusCode int,
	message string,
) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
	}
}