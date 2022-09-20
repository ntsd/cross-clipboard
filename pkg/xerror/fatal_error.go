package xerror

import "fmt"

type FatalError struct {
	Message    string
	WrappedErr error
}

func (e *FatalError) Error() string {
	if e.WrappedErr != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.WrappedErr.Error())
	}
	return e.Message
}

func (e *FatalError) Wrap(err error) *FatalError {
	newErr := *e
	newErr.WrappedErr = err
	return &newErr
}

func (e *FatalError) Unwrap() error {
	return e.WrappedErr
}

func NewFatalError(message string) *FatalError {
	return &FatalError{
		Message: message,
	}
}

func NewFatalErrorf(format string, a ...any) *FatalError {
	return NewFatalError(fmt.Sprintf(format, a...))
}
