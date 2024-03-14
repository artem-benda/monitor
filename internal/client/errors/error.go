package errors

type ErrNetwork struct {
	Err error
}

func (e ErrNetwork) Error() string {
	return e.Err.Error()
}

func (e *ErrNetwork) Unwrap() error {
	return e.Err
}
