// Copyright (c) 2015, Nick Patavalis (npat@efault.net).
// All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.

// +build linux freebsd netbsd openbsd darwin dragonfly
// +build !nocgo

// MakeRaw Termios method implementation that calls LIBC's cfmakeraw()
// through CGo.

package termios

/*
#include <termios.h>
*/
import "C"

// MakeRaw sets the terminal attributes in t to values appropriate for
// configuring the terminal to "raw" mode: Input available character
// by character, echoing disabled, and all special processing of
// terminal input and output characters disabled. Notice that MakeRaw
// does not actually configure the terminal, it only sets the
// attributes in t. In order to configure the terminal, you must
// subsequently call the t.SetFd method.
func (t *Termios) MakeRaw() {
	C.cfmakeraw(&t.t)
}
