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

var (
	// ErrClosed is returned by package functions and methods to
	// indicate failure because the port they tried to operate on
	// has been closed.
	ErrClosed = mkErr(efClosed, "port closed")
	// ErrTimeout is returned by Port methods Read and Write to
	// indicate that the operation took too long, and the set
	// timeout or deadline has expired.
	ErrTimeout = mkErr(efTimeout, "timeout/deadline expired")
	// Returned by Port method Read, in accordance with the io.Reader
	// interface.
	ErrEOF = io.EOF
	// Returned by Port method Write, in accordance with the
	// io.Writer interface.
	ErrUnexpectedEOF = io.ErrUnexpectedEOF
	// Other errors may as well be returned by package methods and
	// functions. These may be system-dependent and subject to
	// change, so you should not act based on their value.
)
