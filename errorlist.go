package errors

// ErrorList represents a list of errors
type ErrorList struct {
	errs []error
}

// Err will return an error value if the ErrorList isn't empty
func (e *ErrorList) Err() (err error) {
	if len(e.errs) == 0 {
		return
	}

	return e
}

// Push will push an error to the error list
func (e *ErrorList) Push(err error) {
	if err == nil {
		return
	}

	e.errs = append(e.errs, err)
}

// Copy will copy a slice of errors to the error list
func (e *ErrorList) Copy(errs []error) {
	for _, err := range errs {
		e.Push(err)
	}
}

// Error will return the error string value
func (e *ErrorList) Error() (out string) {
	switch len(e.errs) {
	case 0:
		return
	case 1:
		return e.errs[0].Error()
	}

	var bs []byte
	for _, err := range e.errs {
		bs = append(bs, []byte(err.Error())...)
		bs = append(bs, ',', '\n')
	}

	return string(bs)
}
