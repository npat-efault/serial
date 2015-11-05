// Copyright (c) 2015, Nick Patavalis (npat@efault.net).
// All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.

// +build linux freebsd netbsd openbsd darwin dragonfly solaris
// +build nocgo !cgo

// Standard termios constants. Should be available for on systems. See
// also "termios_conts_<system>.go" which defines additional constants
// available on some systems, or provides defaults for them (if they
// are not available).

package termios

import "golang.org/x/sys/unix"

const (
	// Input mode flags (access with Termios.IFlag())
	BRKINT = unix.BRKINT
	ICRNL  = unix.ICRNL
	IGNBRK = unix.IGNBRK
	IGNCR  = unix.IGNCR
	IGNPAR = unix.IGNPAR
	INLCR  = unix.INLCR
	INPCK  = unix.INPCK
	ISTRIP = unix.ISTRIP
	IXANY  = unix.IXANY
	IXOFF  = unix.IXOFF
	IXON   = unix.IXON
	PARMRK = unix.PARMRK

	// Output mode flags (access with Termios.OFlag())
	OPOST = unix.OPOST

	// Control mode flags (access with Termios.CFlag())
	CSIZE  = unix.CSIZE
	CS5    = unix.CS5
	CS6    = unix.CS6
	CS7    = unix.CS7
	CS8    = unix.CS8
	CSTOPB = unix.CSTOPB
	CREAD  = unix.CREAD
	PARENB = unix.PARENB
	PARODD = unix.PARODD
	HUPCL  = unix.HUPCL
	CLOCAL = unix.CLOCAL

	// Local mode flags (access with Termios.LFlag())
	ECHO   = unix.ECHO
	ECHOE  = unix.ECHOE
	ECHOK  = unix.ECHOK
	ECHONL = unix.ECHONL
	ICANON = unix.ICANON
	IEXTEN = unix.IEXTEN
	ISIG   = unix.ISIG
	NOFLSH = unix.NOFLSH
	TOSTOP = unix.TOSTOP

	// Cc subscript names (access with Termios.{Cc,CcSet})
	VEOF   = unix.VEOF
	VEOL   = unix.VEOL
	VERASE = unix.VERASE
	VINTR  = unix.VINTR
	VKILL  = unix.VKILL
	VMIN   = unix.VMIN
	VQUIT  = unix.VQUIT
	VSTART = unix.VSTART
	VSTOP  = unix.VSTOP
	VSUSP  = unix.VSUSP
	VTIME  = unix.VTIME

	// Values for the act argument of Termios.SetFd
	TCSANOW   = 0
	TCSADRAIN = 1
	TCSAFLUSH = 2

	// Queue selectors for function Flush
	TCIFLUSH  = unix.TCIFLUSH
	TCOFLUSH  = unix.TCOFLUSH
	TCIOFLUSH = unix.TCIOFLUSH
)
