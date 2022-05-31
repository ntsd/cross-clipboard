package xerror

type FatalError struct {
	err error
}

func (e *FatalError) Error() string { return e.err.Error() }

func NewFatalError(err error) *FatalError {
	return &FatalError{err}
}

type RuntimeError struct {
	err error
}

func (e *RuntimeError) Error() string { return e.err.Error() }

func NewRuntimeError(err error) *RuntimeError {
	return &RuntimeError{err}
}
