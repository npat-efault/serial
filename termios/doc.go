// Copyright (c) 2015, Nick Patavalis (npat@efault.net).
// All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.

// Package termios is a simple Go wrapper to the POSIX Terminal
// Interface (POSIX Termios). See [1].
//
// Similatities to the C API
//
// The API exported by this package is very similar, though not
// identical, to the POSIX Termios C API. Some familiarity with the
// latter is required in order to use this package. See termios(3) or
// [2]. Any differences are covered by the package documentation.
//
// The equivalent to the POSIX termios structure is defined as type
// Termios in this package, which is a structure with *no* exported
// fields. To access the standard POSIX termios flag fields (c_cflag,
// c_oflag, etc) use the methods:
//
//    Termios.CFlag() // Returns ptr to the c_cflag field equivalent.
//    Termios.OFlag() // Returns ptr to the c_oflag field equivalent.
//    ... etc ...
//
// The pointers returned by these methods are of type *TcFlag, wich
// provides convinience methods to set, clear, and test individual
// flags or groups of flags. E.g:
//
//    Termios.CFlag().Clr(termios.CSIZE).Set(termios.CS8)
//    Termios.CFlag().Clr(termios.CLOCAL|termios.HUPCL)
//
// To access the contol-characters vector in termios (c_cc) use the
// following methods.
//
//    Termios.Cc(i) // Return control-char @ index i
//    Termios.CcSet(i, c) // Set control-char @ index i to c.
//
// C API functions that take a termios-structure argument, like
// tc{get,set}attr, cf{set,get}{i,o}speed, etc. are mapped to Termios
// methods:
//
//    tcsetattr --> Termios.SetFd
//    tcgetattr --> Termios.GetFd
//    cfget{i,o}speed --> Termios.Get{I,O}Speed
//    cfset{i,o}speed --> Termios.Set{I,O}Speed
//    cfmakeraw --> Termios.MakeRaw
//
// Unlike their C-API equivalent the speed-setting and speed-getting
// methods take and return numeric baudrate values (integers)
// expressed in bits-per-second (not Bxxx speed-codes).
//
// C API functions that operate directly on fd's are mapped to
// similarly named functions:
//
//    tcflush --> Flush
//    tcdrain --> Drain
//    tcsendbreak --> SendBreak
//
// C API constants (macros), such as mode-flags (e.g. CLOCAL, OPOST),
// retain their names.
//
// Supported systems
//
// Package termios should work on all systems that support the POSIX
// terminal interface, that is, on most Unix-like systems.  Depending
// on the system, package termios can either be built to use the
// system's LIBC functions and macros through CGo, or as a pure-go
// package that issues system calls directly. In either case, the API
// presented by the package is identical.
//
// References
//
// [1] POSIX Terminal Interface
// (https://en.wikipedia.org/wiki/POSIX_terminal_interface)
//
// [2] Single Unix Specification V2, General Terminal interface
// (http://pubs.opengroup.org/onlinepubs/7908799/xbd/termios.html)
package termios
