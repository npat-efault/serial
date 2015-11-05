// Copyright (c) 2015, Nick Patavalis (npat@efault.net).
// All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.

// +build linux freebsd netbsd openbsd darwin dragonfly solaris

package termios

// speedTable is used to map numeric tty speeds (baudrates) to the
// respective system-specific code values. It may or may-not be used
// by system-specific implementations.
type speedTable []struct {
	speed int
	code  spdCode
}

func (t speedTable) Code(speed int) (code spdCode, ok bool) {
	for _, s := range t {
		if s.speed == speed {
			return s.code, true
		}
	}
	return 0, false
}

func (t speedTable) Speed(code spdCode) (speed int, ok bool) {
	for _, s := range t {
		if s.code == code {
			return s.speed, true
		}
	}
	return 0, false
}
