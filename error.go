package errors

const (
	// ErrIsClosed is returned when an action is attempted on a closed service
	ErrIsClosed = Error("cannot perform action, service is closed")
)

// New will return a new error string
func New(str string) (err Error) {
	return Error(str)
}

// Error is a basic error
type Error string

func (e Error) Error() string {
	return string(e)
}
