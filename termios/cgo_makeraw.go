// Copyright (c) 2015, Nick Patavalis (npat@efault.net).
// All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.

// +build linux freebsd netbsd openbsd darwin dragonfly

// MakeRaw Termios method implementation that calls LIBC's cfmakeraw()
// through CGo.

package termios

/*
#include <termios.h>
*/
import "C"

func (t *Termios) MakeRaw() {
	C.cfmakeraw(&t.t)
}
