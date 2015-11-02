// Copyright (c) 2015, Nick Patavalis (npat@efault.net).
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.txt file.

// +build linux freebsd netbsd openbsd darwin dragonfly

// cfMakeRaw function that calls system's LIBC equivalent using CGo.
//
// Most POSIX systems provide a cfmakeraw(3) LIBC function for setting
// a tty to "raw mode". Unfortunatelly, this function is not strictly
// a POSIX requirement, so some systems (e.g. solaris) ommit it. If
// this is the case with your system use file "cfmakeraw_missing.go"
// instead, by editing the build-tags accordingly.

package serial

/*
#include <termios.h>
*/
import "C"

func cfMakeRaw(tios *C.struct_termios) {
	C.cfmakeraw(tios)
}
