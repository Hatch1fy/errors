package errors

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ErrorList represents a list of errors.
type ErrorList struct {
	errs []error
	*stack
}

// Err will return an error value if the ErrorList isn't empty.
// Err also records the stack trace at the point it was called,
// if retuning the whole error list as an error.
func (e *ErrorList) Err() (err error) {
	switch len(e.errs) {
	case 0:
		return nil
	case 1:
		return e.errs[0]
	default:
		e.stack = callers()
		return e
	}
}

// ErrorWithStackTrace is a minimal interface to check that error contains a StackTrace.
type ErrorWithStackTrace interface {
	StackTrace() StackTrace
	Error() string
}

// Push will push an error to the error list. If no message parts are provided,
// the original error will be preserved and no wrapping will happen unless error
// has no stack trace attached. But if message parts are provided, the error will be
// wrapped with a message and a new stack trace. Refer to Cause() to get the original error.
func (e *ErrorList) Push(err error, message ...string) {
	if err == nil {
		return
	}

	if len(message) > 0 {
		// message parts provided, need to wrap the error
		err = &withMessage{
			cause: err,
			msg:   strings.Join(message, " "),
		}
		e.errs = append(e.errs, &withStack{
			error: err,
			stack: callers(),
		})
		return
	}

	switch v := err.(type) {
	case *withStack, *fundamental:
		// append as-is our internal error with stack
		e.errs = append(e.errs, err)
	case *ErrorList:
		// append underlying errors from a list
		e.errs = append(e.errs, v.errs...)
	default:
		if stErr, ok := err.(ErrorWithStackTrace); ok {
			// append the original ErrorWithStackTrace as-is
			e.errs = append(e.errs, stErr)
			return
		}
		// enrich the error with its own stack trace
		e.errs = append(e.errs, &withStack{
			error: err,
			stack: callers(),
		})
	}
}

// Pushf will push an error to the error list with an additional message as per fmt.Sprintf rules.
func (e *ErrorList) Pushf(err error, format string, args ...interface{}) {
	if err == nil {
		return
	}
	err = &withMessage{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}
	e.errs = append(e.errs, &withStack{
		error: err,
		stack: callers(),
	})
}

// Copy will copy a slice of errors to the error list.
func (e *ErrorList) Copy(errs []error) {
	for _, err := range errs {
		if err == nil {
			continue
		}

		switch v := err.(type) {
		case *withStack, *fundamental:
			// append as-is our internal error with stack
			e.errs = append(e.errs, err)
		case *ErrorList:
			// append underlying errors from a list
			e.errs = append(e.errs, v.errs...)
		default:
			if stErr, ok := err.(ErrorWithStackTrace); ok {
				// append the original ErrorWithStackTrace as-is
				e.errs = append(e.errs, stErr)
				continue
			}
			// enrich the error with its own stack trace
			e.errs = append(e.errs, &withStack{
				error: err,
				stack: callers(),
			})
		}
	}
}

// ErrorListIterator allows to iterate errors in the error list.
type ErrorListIterator interface {
	Len() (n int)
	ForEach(fn func(error) (end bool))
	Error() string
}

// Check that *ErrorList implements ErrorListIterator at compile time.
var _ ErrorListIterator = &ErrorList{}

// ForEach will iterate through a list of errors.
func (e *ErrorList) ForEach(fn func(error) (end bool)) {
	for _, err := range e.errs {
		if fn(err) {
			return
		}
	}
}

// Len will return the number of errors within the error list.
func (e *ErrorList) Len() (n int) {
	return len(e.errs)
}

// Error will return the error string value.
func (e *ErrorList) Error() string {
	switch len(e.errs) {
	case 0:
		return ""
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

// MarshalJSON is a json encoding helper func.
func (e ErrorList) MarshalJSON() (bs []byte, err error) {
	errs := make([]string, 0, len(e.errs))
	for _, err := range e.errs {
		errs = append(errs, err.Error())
	}

	return json.Marshal(errs)
}

// UnmarshalJSON is a json decoding helper func.
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
