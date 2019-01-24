package errors

import (
	"encoding/json"
)

// ErrorList represents a list of errors
type ErrorList struct {
	errs []error
}

// Err will return an error value if the ErrorList isn't empty
func (e *ErrorList) Err() (err error) {
	switch len(e.errs) {
	case 0:
		return nil
	case 1:
		return e.errs[0]
	default:
		return e
	}
}

// Push will push an error to the error list
func (e *ErrorList) Push(err error) {
	if err == nil {
		return
	}

	switch v := err.(type) {
	case *ErrorList:
		e.errs = append(e.errs, v.errs...)
	default:
		e.errs = append(e.errs, err)
	}
}

// Copy will copy a slice of errors to the error list
func (e *ErrorList) Copy(errs []error) {
	for _, err := range errs {
		e.Push(err)
	}
}

// ForEach will iterate through a list of errors
func (e *ErrorList) ForEach(fn func(error) (end bool)) {
	for _, err := range e.errs {
		if fn(err) {
			return
		}
	}
}

// Len will return the number of errors within the error list
func (e *ErrorList) Len() (n int) {
	return len(e.errs)
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

// MarshalJSON is a json encoding helper func
func (e ErrorList) MarshalJSON() (bs []byte, err error) {
	errs := make([]string, 0, len(e.errs))
	for _, err := range e.errs {
		errs = append(errs, err.Error())
	}

	return json.Marshal(errs)
}

// UnmarshalJSON is a json decoding helper func
func (e *ErrorList) UnmarshalJSON(bs []byte) (err error) {
	var errs []string
	if err = json.Unmarshal(bs, &errs); err != nil {
		return
	}

	for _, errStr := range errs {
		e.errs = append(e.errs, Error(errStr))
	}

	return
}
