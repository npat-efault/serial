// Copyright (c) 2015, Nick Patavalis (npat@efault.net).
// All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.

// +build linux
// +build nocgo !cgo

// Pure-Go implementation for linux that issues system-calls directly

package termios

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

// TcFlag is a word containing Termios mode-flags
type TcFlag uint32

// Cc is a Termios control character
type Cc uint8

// spdCode is a Termios speed code
type spdCode uint32

// Termios is the terminal-attributes structure
type Termios struct {
	t unix.Termios
}

// IFlag returns a pointer to the Termios field keeping the input-mode
// flags
func (t *Termios) IFlag() *TcFlag { return (*TcFlag)(&t.t.Iflag) }

// OFlag returns a pointer to the Termios field keeping the output-mode
// flags
func (t *Termios) OFlag() *TcFlag { return (*TcFlag)(&t.t.Oflag) }

// LFlag returns a pointer to the Termios field keeping the local-mode
// flags
func (t *Termios) LFlag() *TcFlag { return (*TcFlag)(&t.t.Lflag) }

// CFlag returns a pointer to the Termios field keeping the control-mode
// flags
func (t *Termios) CFlag() *TcFlag { return (*TcFlag)(&t.t.Cflag) }

// Cc returns the Termios control character in index idx
func (t *Termios) Cc(idx int) Cc { return Cc(t.t.Cc[idx]) }

// CcSet sets the Termios control character in index idx to c.
func (t *Termios) CcSet(idx int, c Cc) { t.t.Cc[idx] = uint8(c) }

// ioctlV performs an ioctl with a value (integer) input argument
func ioctlV(fd int, req int, v uintptr) error {
	_, _, err := unix.Syscall(unix.SYS_IOCTL,
		uintptr(fd), uintptr(req), v)
	if err != 0 {
		return err
	}
	return nil
}

// ioctlP performs an ioctl with a pointer input or output argument
func ioctlP(fd int, req int, p unsafe.Pointer) error {
	_, _, err := unix.Syscall(unix.SYS_IOCTL,
		uintptr(fd), uintptr(req), uintptr(p))
	// Not strictly required, but better be safe.
	use(p)
	if err != 0 {
		return err
	}
	return nil
}

// SetFd configures the terminal corresponding to the file-descriptor
// fd using the attributes in the Termios structure t. Argument act
// must be one of the constants TCSANOW, TCSADRAIN, TCSAFLUSH. See
// tcsetattr(3) for more. If the system call is interrupted by a
// signal, SetFd retries it automatically.
func (t *Termios) SetFd(fd int, act int) error {
	var req int
	switch act {
	case TCSANOW:
		req = unix.TCSETS
	case TCSADRAIN:
		req = unix.TCSETSW
	case TCSAFLUSH:
		req = unix.TCSETSF
	default:
		return syscall.EINVAL
	}
	for {
		err := ioctlP(fd, req, unsafe.Pointer(&t.t))
		if err == syscall.EINTR {
			continue
		}
		return err
	}
}

// GetFd reads the attributes of the terminal corresponding to the
// file-descriptor fd and stores them in the Termios structure t. See
// tcgetattr(3) for more.
func (t *Termios) GetFd(fd int) error {
	for {
		err := ioctlP(fd, unix.TCGETS, unsafe.Pointer(&t.t))
		// This is most-likely not possible, but
		// better be safe.
		if err == syscall.EINTR {
			continue
		}
		return err
	}
}

// SetOSpeed sets the output (transmitter) baudrate in termios
// structure t to speed. Argument speed must be a numerical (integer)
// baudrate value in bits-per-second. Returns syscall.EINVAL if the
// requested baudrate is not supported. See also cfsetospeed(3).
func (t *Termios) SetOSpeed(speed int) error {
	code, ok := stdSpeeds.Code(speed)
	if !ok {
		return syscall.EINVAL
	}
	// Standard speed
	t.CFlag().Clr(unix.CBAUD | unix.CBAUDEX).Set(TcFlag(code))
	return nil
}

// SetISpeed sets the input (receiver) baudrate in Termios structure t
// to speed. Argument speed must be a numerical (integer) baudrate
// value in bits-per-second. Returns syscall.EINVAL if the requested
// baudrate is not supported. See also cfsetispeed(3).
func (t *Termios) SetISpeed(speed int) error {
	code, ok := stdSpeeds.Code(speed)
	if !ok {
		return syscall.EINVAL
	}
	// Standard speed
	t.CFlag().Clr((unix.CBAUD | unix.CBAUDEX) << unix.IBSHIFT)
	t.CFlag().Set(TcFlag(code) << unix.IBSHIFT)
	return nil
}

// GetOSpeed returns the output (transmitter) baudrate in Termios
// structure t as a numerical (integer) value in
// bits-per-second. Returns err == syscal.EINVAL if the baudrate in t
// cannot be decoded. See also getospeed(3).
func (t *Termios) GetOSpeed() (speed int, err error) {
	c := t.CFlag().Msk(unix.CBAUD | unix.CBAUDEX)
	if c != unix.BOTHER {
		// Standard speed
		speed, ok := stdSpeeds.Speed(spdCode(c))
		if !ok {
			return 0, syscall.EINVAL
		}
		return speed, nil
	} else {
		// Custom speed
		return 0, nil
	}
}

// GetISpeed returns the input (receiver) baudrate in Termios
// structure t as a numerical (integer) value in
// bits-per-second. Returns err == syscal.EINVAL if the baudrate in t
// cannot be decoded. See also getispeed(3).
func (t *Termios) GetISpeed() (speed int, err error) {
	c := t.CFlag().Msk((unix.CBAUD | unix.CBAUDEX) << unix.IBSHIFT)
	c >>= unix.IBSHIFT
	if c != unix.BOTHER {
		// Standard speed
		speed, ok := stdSpeeds.Speed(spdCode(c))
		if !ok {
			return 0, syscall.EINVAL
		}
		return speed, nil
	} else {
		// Custom speed
		return 0, nil
	}
}

// Flush discards data received but not yet read (input queue), and/or
// data written but not yet transmitted (output queue), depending on
// the value of the qsel argument. Argument qsel must be one of the
// constants TCIFLUSH (flush input queue), TCOFLUSH (flush output
// queue), TCIOFLUSH (flush both queues). See also tcflush(3).
func Flush(fd int, qsel int) error {
	for {
		err := ioctlV(fd, unix.TCFLSH, uintptr(qsel))
		// This is most-likely not possible, but
		// better be safe.
		if err == syscall.EINTR {
			continue
		}
		return err
	}
}

// Drain blocks until all data written to the terminal fd are
// transmitted. See also tcdrain(3). If the system call is interrupted
// by a signal, Drain retries it automatically.
func Drain(fd int) error {
	for {
		err := ioctlV(fd, unix.TCSBRK, 1)
		if err == syscall.EINTR {
			continue
		}
		return err
	}
}

// SendBreak sends a continuous stream of zero bits to the terminal
// corresponding to file-descriptor fd, lasting between 0.25 and 0.5
// seconds.
func SendBreak(fd int) error {
	return ioctlV(fd, unix.TCSBRK, 0)
}

// Flow suspends or resumes the transmission or reception of data on
// the terminal associated with fd, depending on the value of the act
// argument. The act argument value must be one of: TCOOFF (suspend
// transmission), TCOON (resume transmission), TCIOFF (suspend
// reception by sending a STOP char), TCION (resume reception by
// sending a START char). See also tcflow(3).
func Flow(fd int, act int) error {
	for {
		err := ioctlV(fd, unix.TCXONC, uintptr(act))
		// This is most-likely not possible, but
		// better be safe.
		if err == syscall.EINTR {
			continue
		}
		return err
	}
}

// GetPgrp returns the process group ID of the foreground process
// group associated with the terminal. See tcgetpgrp(3).
func GetPgrp(fd int) (pgid int, err error) {
	err = ioctlP(fd, unix.TIOCGPGRP, unsafe.Pointer(&pgid))
	if err != nil {
		return 0, err
	}
	return pgid, nil
}

// SetPgrp sets the foreground process group ID associated with the
// terminal to pgid. See tcsetpgrp(3).
func SetPgrp(fd int, pgid int) error {
	return ioctlP(fd, unix.TIOCSPGRP, unsafe.Pointer(&pgid))
}
