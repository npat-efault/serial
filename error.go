package serial

// ErrCode encodes the errors returned by this package.
type ErrCode int

const (
	ErrUnknown ErrCode = iota
	ErrOpen            // Cannot open port, see system error
	ErrConf            // Cannot set port configuration
	ErrGetConf         // Cannot read / decode port configuration
	ErrRead            // Read from port failed
	ErrWrite           // Write to port failed
	ErrTimeout         // Read or write operation timed-out
	ErrReset           // Cannot reset port to original settings
	ErrClosed          // Port has been closed
	ErrInvalid         // Invalid argument or parameter
)

var errStr = [...]string{
	"unkown error",
	"cannot open port",
	"cannot config port",
	"cannot read conf",
	"read failed",
	"write failed",
	"timeout expired",
	"cannot reset port",
	"port closed",
	"invalid argument",
}

// Error is the error type returned by all functions and methods in
// this package. The Code field is the system-independent error code
// filled-in by the package (see ErCode constants). The field Err
// supplies the system-specific error, if available. In most cases it
// should be enough to examine in order to act upon the erroe.
type Error struct {
	Code ErrCode
	Err  error // system error
}

func (e *Error) Error() string {
	if e.Err == nil {
		return errStr[e.Code]
	} else {
		return errStr[e.Code] + ": " + e.Err.Error()
	}
}

func (e *Error) Timeout() bool {
	return e.Code == ErrTimeout
}

func (e *Error) Temporary() bool {
	return e.Code == ErrTimeout
}

func (e *Error) Closed() bool {
	return e.Code == ErrClosed
}
