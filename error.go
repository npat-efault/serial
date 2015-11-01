package serial

import "io"

type errFlags int

const (
	efClosed errFlags = 1 << iota
	efTimeout
	efTemporary
)

type errT struct {
	flags errFlags
	msg   string
}

func (e *errT) Error() string {
	return e.msg
}

func (e *errT) Timeout() bool {
	return e.flags&efTimeout != 0
}

func (e *errT) Temporary() bool {
	return e.flags&(efTimeout|efTemporary) != 0
}

func (e *errT) Closed() bool {
	return e.flags&efClosed != 0
}

func mkErr(flags errFlags, msg string) error {
	return &errT{flags: flags, msg: msg}
}

func newErr(msg string) error {
	return &errT{msg: msg}
}

// Errors returned by functions and methods in this package. Other
// errors may as well be returned which are system-dependent and
// subject to change. You should not act based on the value of
// errors other than these.
var (
	// ErrClosed is returned by package functions and methods to
	// indicate failure because the port they tried to operate on
	// has been closed. ErrClosed has Closed() == true.
	ErrClosed = mkErr(efClosed, "port closed")
	// ErrTimeout is returned by Port methods Read and Write to
	// indicate that the operation took too long, and the set
	// timeout or deadline has expired. ErrTimeout has Timeout()
	// == true and Temporary() == true
	ErrTimeout = mkErr(efTimeout, "timeout/deadline expired")
	// ErrEOF is returned by Port method Read, in accordance with
	// the io.Reader interface.
	ErrEOF = io.EOF
	// ErrUnexpectedEOF is returned by Port method Write, in
	// accordance with the io.Writer interface.
	ErrUnexpectedEOF = io.ErrUnexpectedEOF
)
