// Copyright (c) 2015, Nick Patavalis (npat@efault.net).
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.txt file.

// +build linux
// +build !nocgo

// Standard serial port speeds. Linux set (externds POSIX). Taken from
// system's LIBC headers using CGo.

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
	{57600, C.B57600},
	{115200, C.B115200},
	{230400, C.B230400},
	{460800, C.B460800},
	{500000, C.B500000},
	{576000, C.B576000},
	{921600, C.B921600},
	{1000000, C.B1000000},
	{1152000, C.B1152000},
	{2000000, C.B2000000},
	{2500000, C.B2500000},
	{3000000, C.B3000000},
	{3500000, C.B3500000},
	{4000000, C.B4000000},
}
