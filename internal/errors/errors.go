package errors

type DefaultError struct {
	Code    int
	Message string
}

type NotExists struct {
	DefaultError
}

func (e NotExists) Error() string {
	return e.Message
}

func NewNotExists(message string) error {
	return &NotExists{
		DefaultError{
			Code:    100,
			Message: message,
		},
	}
}

type UnexpectedResult struct {
	DefaultError
}

func (e UnexpectedResult) Error() string {
	return e.Message
}

func NewUnexpectedResult(message string) error {
	return &UnexpectedResult{
		DefaultError{
			Code:    101,
			Message: message,
		},
	}
}
