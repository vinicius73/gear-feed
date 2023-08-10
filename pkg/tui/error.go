package tui

type MsgError struct {
	err error
}

func Error(err error) MsgError {
	return MsgError{err}
}

func (e MsgError) Error() string {
	return e.err.Error()
}

func (e MsgError) Unwrap() error {
	return e.err
}
