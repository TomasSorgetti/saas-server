package errors

type HTTPError struct {
	Code    int
	Message string
	Details any
}

func (e *HTTPError) Error() string {
	return e.Message
}

func New(code int, message string, details ...any) *HTTPError {
	var d any
	if len(details) > 0 {
		d = details[0]
	}
	return &HTTPError{Code: code, Message: message, Details: d}
}