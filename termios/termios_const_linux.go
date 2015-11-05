// Copyright (c) 2015, Nick Patavalis (npat@efault.net).
// All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.

// +build linux
// +build nocgo !cgo

// Non-standard (non-POSIX) termios constants for linux.

package termios

import "golang.org/x/sys/unix"

// Non-standard (non-POSIX) termios mode-flags and other constants.
const (
	// Input mode flags (access with Termios.IFlag())
	// Non-standard (0 if unavailable)
	IMAXBEL = unix.IMAXBEL
	IUCLC   = unix.IUCLC
	IUTF8   = unix.IUTF8

	// Control mode flags (access with Termios.CFlag())
	// Non-standard (0 if unavailable)
	CRTSCTS = unix.CRTSCTS
	CMSPAR  = unix.CMSPAR

	// Local mode flags (access with Termios.LFlag())
	// Non-standard (0 if unavailable)
	PENDIN  = unix.PENDIN
	ECHOCTL = unix.ECHOCTL
	ECHOPRT = unix.ECHOPRT
	ECHOKE  = unix.ECHOKE
	FLUSHO  = unix.FLUSHO
	EXTPROC = unix.EXTPROC

	// Cc subscript names (access with Termios.{Cc,CcSet})
	// Non-standard (-1 if unavailable)
	VREPRINT = unix.VREPRINT
	VDISCARD = unix.VDISCARD
	VWERASE  = unix.VWERASE
	VLNEXT   = unix.VLNEXT
	VEOL2    = unix.VEOL2
)
