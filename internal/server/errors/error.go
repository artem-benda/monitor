package errors

type ErrStorageConnection struct {
	Err error
}

func (e ErrStorageConnection) Error() string {
	return e.Err.Error()
}

func (e *ErrStorageConnection) Unwrap() error {
	return e.Err
}
