// Copyright (c) 2015, Nick Patavalis (npat@efault.net).
// All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.

// +build linux freebsd netbsd openbsd darwin dragonfly solaris
// +build !nocgo

// Implementation that uses CGo to access the system's LIBC functions
// and macros.

package termios

/*
#include <termios.h>
#include <unistd.h>

#include "termios_nonstd.h"
*/
import "C"
import "syscall"

const (
	// Input mode flags (access with Termios.IFlag())
	BRKINT = C.BRKINT
	ICRNL  = C.ICRNL
	IGNBRK = C.IGNBRK
	IGNCR  = C.IGNCR
	IGNPAR = C.IGNPAR
	INLCR  = C.INLCR
	INPCK  = C.INPCK
	ISTRIP = C.ISTRIP
	IXANY  = C.IXANY
	IXOFF  = C.IXOFF
	IXON   = C.IXON
	PARMRK = C.PARMRK
	// Non-standard (0 if unavailable)
	IMAXBEL = C.IMAXBEL
	IUCLC   = C.IUCLC
	IUTF8   = C.IUTF8

	// Output mode flags (access with Termios.OFlag())
	OPOST = C.OPOST

	// Control mode flags (access with Termios.CFlag())
	CSIZE  = C.CSIZE
	CS5    = C.CS5
	CS6    = C.CS6
	CS7    = C.CS7
	CS8    = C.CS8
	CSTOPB = C.CSTOPB
	CREAD  = C.CREAD
	PARENB = C.PARENB
	PARODD = C.PARODD
	HUPCL  = C.HUPCL
	CLOCAL = C.CLOCAL
	// Non-standard (0 if unavailable)
	CRTSCTS = C.CRTSCTS
	CMSPAR  = C.CMSPAR

	// Local mode flags (access with Termios.LFlag())
	ECHO   = C.ECHO
	ECHOE  = C.ECHOE
	ECHOK  = C.ECHOK
	ECHONL = C.ECHONL
	ICANON = C.ICANON
	IEXTEN = C.IEXTEN
	ISIG   = C.ISIG
	NOFLSH = C.NOFLSH
	TOSTOP = C.TOSTOP
	// Non-standard (0 if unavailable)
	PENDIN  = C.PENDIN
	ECHOCTL = C.ECHOCTL
	ECHOPRT = C.ECHOPRT
	ECHOKE  = C.ECHOKE
	FLUSHO  = C.FLUSHO
	EXTPROC = C.EXTPROC

	// Cc subscript names (access with Termios.{Cc,CcSet})
	VEOF   = C.VEOF
	VEOL   = C.VEOL
	VERASE = C.VERASE
	VINTR  = C.VINTR
	VKILL  = C.VKILL
	VMIN   = C.VMIN
	VQUIT  = C.VQUIT
	VSTART = C.VSTART
	VSTOP  = C.VSTOP
	VSUSP  = C.VSUSP
	VTIME  = C.VTIME
	// Non-standard (-1 if unavailable)
	VREPRINT = C.VREPRINT
	VDISCARD = C.VDISCARD
	VWERASE  = C.VWERASE
	VLNEXT   = C.VLNEXT
	VEOL2    = C.VEOL2

	// Values for the act argument of Termios.SetFd
	TCSANOW   = C.TCSANOW
	TCSADRAIN = C.TCSADRAIN
	TCSAFLUSH = C.TCSAFLUSH

	// Queue selectors for function Flush
	TCIFLUSH  = C.TCIFLUSH
	TCOFLUSH  = C.TCOFLUSH
	TCIOFLUSH = C.TCIOFLUSH

	// Values for the act argument of Flow
	TCOOFF = C.TCOOFF
	TCOON  = C.TCOON
	TCIOFF = C.TCIOFF
	TCION  = C.TCION
)

// TcFlag is a word containing Termios mode-flags
type TcFlag C.tcflag_t

// Cc is a Termios control character
type Cc C.cc_t

// spdCode is a Termios speed code
type spdCode C.speed_t

// Termios is the terminal-attributes structure
type Termios struct {
	t C.struct_termios
}

// IFlag returns a pointer to the Termios field keeping the input-mode
// flags
func (t *Termios) IFlag() *TcFlag { return (*TcFlag)(&t.t.c_iflag) }

// OFlag returns a pointer to the Termios field keeping the output-mode
// flags
func (t *Termios) OFlag() *TcFlag { return (*TcFlag)(&t.t.c_oflag) }

// LFlag returns a pointer to the Termios field keeping the local-mode
// flags
func (t *Termios) LFlag() *TcFlag { return (*TcFlag)(&t.t.c_lflag) }

// CFlag returns a pointer to the Termios field keeping the control-mode
// flags
func (t *Termios) CFlag() *TcFlag { return (*TcFlag)(&t.t.c_cflag) }

// Cc returns the Termios control character in index idx
func (t *Termios) Cc(idx int) Cc { return Cc(t.t.c_cc[idx]) }

// CcSet sets the Termios control character in index idx to c.
func (t *Termios) CcSet(idx int, c Cc) { t.t.c_cc[idx] = C.cc_t(c) }

// SetFd configures the terminal corresponding to the file-descriptor
// fd using the attributes in the Termios structure t. Argument act
// must be one of the constants TCSANOW, TCSADRAIN, TCSAFLUSH. See
// tcsetattr(3) for more. If the system call is interrupted by a
// signal, SetFd retries it automatically.
func (t *Termios) SetFd(fd int, act int) error {
	for {
		r, err := C.tcsetattr(C.int(fd), C.int(act), &t.t)
		if r < 0 {
			if err == syscall.EINTR {
				continue
			}
			return err
		}
		return nil
	}
}

// GetFd reads the attributes of the terminal corresponding to the
// file-descriptor fd and stores them in the Termios structure t. See
// tcgetattr(3) for more.
func (t *Termios) GetFd(fd int) error {
	for {
		r, err := C.tcgetattr(C.int(fd), &t.t)
		if r < 0 {
			// This is most-likely not possible, but
			// better be safe.
			if err == syscall.EINTR {
				continue
			}
			return err
		}
		return nil
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
	C.cfsetospeed(&t.t, C.speed_t(code))
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
	C.cfsetispeed(&t.t, C.speed_t(code))
	return nil
}

// GetOSpeed returns the output (transmitter) baudrate in Termios
// structure t as a numerical (integer) value in
// bits-per-second. Returns err == syscal.EINVAL if the baudrate in t
// cannot be decoded. See also getospeed(3).
func (t *Termios) GetOSpeed() (speed int, err error) {
	code := C.cfgetospeed(&t.t)
	speed, ok := stdSpeeds.Speed(spdCode(code))
	if !ok {
		return 0, syscall.EINVAL
	}
	return speed, nil
}

// GetISpeed returns the input (receiver) baudrate in Termios
// structure t as a numerical (integer) value in
// bits-per-second. Returns err == syscal.EINVAL if the baudrate in t
// cannot be decoded. See also getispeed(3).
func (t *Termios) GetISpeed() (speed int, err error) {
	code := C.cfgetispeed(&t.t)
	speed, ok := stdSpeeds.Speed(spdCode(code))
	if !ok {
		return 0, syscall.EINVAL
	}
	return speed, nil
}

// Flush discards data received but not yet read (input queue), and/or
// data written but not yet transmitted (output queue), depending on
// the value of the qsel argument. Argument qsel must be one of the
// constants TCIFLUSH (flush input queue), TCOFLUSH (flush output
// queue), TCIOFLUSH (flush both queues). See also tcflush(3).
func Flush(fd int, qsel int) error {
	for {
		r, err := C.tcflush(C.int(fd), C.int(qsel))
		if r < 0 {
			// This is most-likely not possible, but
			// better be safe.
			if err == syscall.EINTR {
				continue
			}
			return err
		}
		return nil
	}
}

// Drain blocks until all data written to the terminal fd are
// transmitted. See also tcdrain(3). If the system call is interrupted
// by a signal, Drain retries it automatically.
func Drain(fd int) error {
	for {
		r, err := C.tcdrain(C.int(fd))
		if r < 0 {
			if err == syscall.EINTR {
				continue
			}
			return err
		}
		return nil
	}
}

// SendBreak sends a continuous stream of zero bits to the terminal
// corresponding to file-descriptor fd, lasting between 0.25 and 0.5
// seconds.
func SendBreak(fd int) error {
	r, err := C.tcsendbreak(C.int(fd), 0)
	if r < 0 {
		return err
	}
	return nil
}

// Flow suspends or resumes the transmission or reception of data on
// the terminal associated with fd, depending on the value of the act
// argument. The act argument value must be one of: TCOOFF (suspend
// transmission), TCOON (resume transmission), TCIOFF (suspend
// reception by sending a STOP char), TCION (resume reception by
// sending a START char). See also tcflow(3).
func Flow(fd int, act int) error {
	for {
		r, err := C.tcflow(C.int(fd), C.int(act))
		if r < 0 {
			// This is most-likely not possible, but
			// better be safe.
			if err == syscall.EINTR {
				continue
			}
			return err
		}
		return nil
	}
}

// GetPgrp returns the process group ID of the foreground process
// group associated with the terminal. See tcgetpgrp(3).
func GetPgrp(fd int) (pgid int, err error) {
	r, err := C.tcgetpgrp(C.int(fd))
	if r < 0 {
		return 0, err
	}
	return int(r), nil
}

// SetPgrp sets the foreground process group ID associated with the
// terminal to pgid. See tcsetpgrp(3).
func SetPgrp(fd int, pgid int) error {
	r, err := C.mytcsetpgrp(C.int(fd), C.pid_t(pgid))
	if r < 0 {
		return err
	}
	return nil
}
