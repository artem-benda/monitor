package errors

// ErrStorageConnection - Обобщенная ошибка доступа к хранилищу
type ErrStorageConnection struct {
	Err error
}

func (e ErrStorageConnection) Error() string {
	return "temporary storage connection error"
}

func (e *ErrStorageConnection) Unwrap() error {
	return e.Err
}
