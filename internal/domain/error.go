package domain

type LineError struct {
	Err  error
	Line string
}

func (e *LineError) Error() string {
	return e.Err.Error()
}

func NewLineError(err error, line string) *LineError {
	return &LineError{
		Err:  err,
		Line: line,
	}
}
