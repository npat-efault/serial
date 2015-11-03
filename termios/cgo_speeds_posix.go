// Copyright (c) 2015, Nick Patavalis (npat@efault.net).
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.txt file.

// +build freebsd netbsd openbsd darwin dragonfly solaris

// Standard serial port speeds defined by POSIX and taken from
// system's LIBC headers using CGo.
//
// These speeds should be available on all POSIX systems. Most likely
// your system will support more (higher) speeds. If this is the case,
// copy this file as "speeds_<system>.go", add the additional
// supported speeds and edit the build-tags accordingly.

package termios

/*
#include <termios.h>
*/
import "C"

var stdSpeeds = speedTable{
	{0, C.B0},
	{50, C.B50},
	{75, C.B75},
	{110, C.B110},
	{134, C.B134},
	{150, C.B150},
	{200, C.B200},
	{300, C.B300},
	{600, C.B600},
	{1200, C.B1200},
	{1800, C.B1800},
	{2400, C.B2400},
	{4800, C.B4800},
	{9600, C.B9600},
	{19200, C.B19200},
	{38400, C.B38400},
}
