package tui

type ErrorMsg struct {
	err error
}

func Error(err error) ErrorMsg {
	return ErrorMsg{err}
}

func (e ErrorMsg) Error() string {
	return e.err.Error()
}

func (e ErrorMsg) Unwrap() error {
	return e.err
}
