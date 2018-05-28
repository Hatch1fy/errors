package errors

// ErrorList represents a list of errors
type ErrorList []error

// Err will return an error value if the ErrorList isn't empty
func (e ErrorList) Err() (err error) {
	if len(e) == 0 {
		return
	}

	return e
}

func (e ErrorList) Error() (out string) {
	switch len(e) {
	case 0:
		return
	case 1:
		return e[0].Error()
	}

	var bs []byte
	for _, err := range e {
		bs = append(bs, []byte(err.Error())...)
		bs = append(bs, ',', '\n')
	}

	return string(bs)
}
