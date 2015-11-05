// Copyright (c) 2015, Nick Patavalis (npat@efault.net).
// All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.

// +build linux freebsd netbsd openbsd darwin dragonfly solaris
// +build nocgo !cgo solaris

// MakeRaw Termios method implementation for pure-Go builds and for
// systems that don't have a cfmakeraw in LIBC.

package termios

func (t *Termios) MakeRaw() {
	t.IFlag().Clr(IMAXBEL | IXOFF | INPCK | BRKINT | PARMRK |
		ISTRIP | INLCR | IGNCR | ICRNL | IXON | IGNPAR)
	t.IFlag().Set(IGNBRK)

	t.OFlag().Clr(OPOST)

	t.LFlag().Clr(ECHO | ECHOE | ECHOK | ECHONL | ICANON | ISIG |
		IEXTEN | NOFLSH | TOSTOP | PENDIN)

	t.CFlag().Clr(CSIZE | PARENB)
	t.CFlag().Set(CS8 | CREAD)
	t.CcSet(VMIN, 1)
	t.CcSet(VTIME, 0)
}
