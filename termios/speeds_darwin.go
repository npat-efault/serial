// Copyright (c) 2020, Alexey McSakoff (mcsakoff@gmail.com).
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.txt file.

// +build darwin
// +build nocgo !cgo

// Standard serial port speeds. Darwin set (extends POSIX).

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
	{7200, unix.B7200},
	{9600, unix.B9600},
	{14400, unix.B14400},
	{19200, unix.B19200},
	{28800, unix.B28800},
	{38400, unix.B38400},
	{57600, unix.B57600},
	{76800, unix.B76800},
	{115200, unix.B115200},
	{230400, unix.B230400},
}
