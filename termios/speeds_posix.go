// Copyright (c) 2015, Nick Patavalis (npat@efault.net).
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.txt file.

// +build freebsd netbsd openbsd darwin dragonfly solaris
// +build nocgo !cgo

// Standard serial port speeds defined by POSIX.
//
// These speeds should be available on all POSIX systems. Most likely
// your system will support more (higher) speeds. If this is the case,
// see "speeds_full.go" or copy this file as "speeds_<system>.go", add
// the additional supported speeds and edit the build-tags
// accordingly.

package termios

import "golang.org/x/sys/unix"

var stdSpeeds = speedTable{
	{0, unix.B0},
	{50, unix.B50},
	{75, unix.B75},
	{110, unix.B110},
	{134, unix.B134},
	{150, unix.B150},
	{200, unix.B200},
	{300, unix.B300},
	{600, unix.B600},
	{1200, unix.B1200},
	{1800, unix.B1800},
	{2400, unix.B2400},
	{4800, unix.B4800},
	{9600, unix.B9600},
	{19200, unix.B19200},
	{38400, unix.B38400},
}
