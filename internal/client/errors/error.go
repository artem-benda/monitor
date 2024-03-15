package errors

type ErrNetwork struct {
	Err error
}

func (e ErrNetwork) Error() string {
	return "temporary network error"
}

func (e *ErrNetwork) Unwrap() error {
	return e.Err
}

type ErrServerTemporary struct {
	Err error
}

func (e ErrServerTemporary) Error() string {
	return "temporary server error"
}

func (e *ErrServerTemporary) Unwrap() error {
	return e.Err
}
