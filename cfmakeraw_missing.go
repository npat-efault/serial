// Copyright (c) 2015, Nick Patavalis (npat@efault.net).
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.txt file.

// +build linux

// cfMakeRaw function implemented by directly settings
// C.struct_termios flags, for systems that don't provide a
// C.cfmakeraw.

package serial

/*
#include <termios.h>
#include <unistd.h>
*/
import "C"

func cfMakeRaw(tios *C.struct_termios) {
	tios.c_iflag &^= C.IMAXBEL | C.IXOFF | C.INPCK |
		C.BRKINT | C.PARMRK | C.ISTRIP | C.INLCR |
		C.IGNCR | C.ICRNL | C.IXON | C.IGNPAR
	tios.c_iflag |= C.IGNBRK
	tios.c_oflag &^= C.OPOST
	tios.c_lflag &^= C.ECHO | C.ECHOE | C.ECHOK | C.ECHONL |
		C.ICANON | C.ISIG | C.IEXTEN | C.NOFLSH |
		C.TOSTOP | C.PENDIN
	tios.c_cflag &^= C.CSIZE | C.PARENB
	tios.c_cflag |= C.CS8 | C.CREAD
	tios.c_cc[C.VMIN] = 1
	tios.c_cc[C.VTIME] = 0
}
