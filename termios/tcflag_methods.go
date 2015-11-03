// Copyright (c) 2015, Nick Patavalis (npat@efault.net).
// All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.

// Bit-twiddling convenience methods for type TcFlag

package termios

// Clr clears (lowers) in w the flags that are set in f. It does not
// affect the other flags in w. Returns (a pointer to) w.
func (w *TcFlag) Clr(f TcFlag) *TcFlag { *w &^= f; return w }

// Set sets (raises) in w the flags that are set in f. It does not
// affect the other flags in w. Returns (a pointer to) w.
func (w *TcFlag) Set(f TcFlag) *TcFlag { *w |= f; return w }

// Any returns true if any of flags set in f are also set in w.
func (w *TcFlag) Any(f TcFlag) bool { return *w&f != 0 }

// All returns true if all of the flags set in f are also set in w.
func (w *TcFlag) All(f TcFlag) bool { return *w&f == f }

// Msk returns the flags in w masked by f (that is, the values of the
// flags in w for which the equivalent flags in f are set). It does
// not modify w.
func (w *TcFlag) Msk(f TcFlag) TcFlag { return *w & f }

// Val returns the flags in w (i.e the value pointed-to by w)
func (w *TcFlag) Val() TcFlag { return *w }
