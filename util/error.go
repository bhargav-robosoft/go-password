package util

type CustomError struct {
	Message string
	Status  int
}

func (e *CustomError) Error() string {
	return e.Message
}
