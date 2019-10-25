package errors

const (
	// ErrIsClosed is returned when an action is attempted on a closed service.
	ErrIsClosed = Error("cannot perform action, service is closed")
)

// Error is a basic error. Allows to use constant
// strings as errors without stack trace.
type Error string

// Error implements error interface.
func (e Error) Error() string {
	return string(e)
}
