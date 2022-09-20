package xerror

import "fmt"

type RuntimeError struct {
	Message    string
	WrappedErr error
}

func (e *RuntimeError) Error() string {
	if e.WrappedErr != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.WrappedErr.Error())
	}
	return e.Message
}

func (e *RuntimeError) Wrap(err error) *RuntimeError {
	newErr := *e
	newErr.WrappedErr = err
	return &newErr
}

func (e *RuntimeError) Unwrap() error {
	return e.WrappedErr
}

func NewRuntimeError(message string) *RuntimeError {
	return &RuntimeError{
		Message: message,
	}
}

func NewRuntimeErrorf(format string, a ...any) *RuntimeError {
	return NewRuntimeError(fmt.Sprintf(format, a...))
}
