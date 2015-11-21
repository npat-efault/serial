// Copyright (c) 2015, Nick Patavalis (npat@efault.net).
// All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.

// +build linux freebsd netbsd openbsd darwin dragonfly solaris

// !! Tests are necessarily superficial !!

package termios_test

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/npat-efault/serial/termios"
)

func TestFlags(t *testing.T) {
	ti := termios.Termios{}
	if *ti.IFlag() != 0 || *ti.OFlag() != 0 ||
		*ti.CFlag() != 0 || *ti.LFlag() != 0 {
		t.Fatalf("Termios struct not zero: %+v", ti)
	}
	ti.CFlag().Clr(termios.CSIZE).Set(termios.CS8)
	if ti.CFlag().Msk(termios.CSIZE) != termios.CS8 {
		t.Fatalf("Bad CFlag: %b", *ti.CFlag())
	}
	if !ti.CFlag().Set(termios.CLOCAL).Any(termios.CLOCAL | termios.HUPCL) {
		t.Fatalf("Bad CFlag: %b", ti.CFlag().Val())
	}
	if ti.CFlag().All(termios.CLOCAL | termios.HUPCL) {
		t.Fatalf("Bad CFlag: %b", ti.CFlag().Val())
	}
}

func TestMakeRaw(t *testing.T) {
	ti := termios.Termios{}
	ti.MakeRaw()
	// Check some arbitrary fields
	if ti.CFlag().Msk(termios.CSIZE) != termios.CS8 {
		t.Fatalf("Bad CFlag: %b", ti.CFlag().Val())
	}
	/*
		if !ti.CFlag().All(termios.CREAD) {
			t.Fatalf("Bad CFlag: %b", ti.CFlag().Val())
		}
	*/
	if ti.Cc(termios.VMIN) != 1 {
		t.Fatalf("Bad VMIN: %b", ti.Cc(termios.VMIN))
	}
	if ti.Cc(termios.VTIME) != 0 {
		t.Fatalf("Bad VTIME: %b", ti.Cc(termios.VTIME))
	}
}

func TestSpeed(t *testing.T) {
	ti := termios.Termios{}
	spds_ok := []int{1200, 4800, 9600, 19200, 38400}
	spds_fail := []int{-1, -2, -3}
	for _, s := range spds_ok {
		// Output Speeds
		if err := ti.SetOSpeed(s); err != nil {
			t.Fatalf("Cannot set out speed %d: %v", s, err)
		}
		s1, err := ti.GetOSpeed()
		if err != nil {
			t.Fatalf("Cannot get out speed %d: %v", s, err)
		}
		if s1 != s {
			t.Fatalf("Out speed: %d != %d", s1, s)
		}
		// Input Speeds
		if err := ti.SetISpeed(s); err != nil {
			t.Fatalf("Cannot set in speed %d: %v", s, err)
		}
		s1, err = ti.GetISpeed()
		if err != nil {
			t.Fatalf("Cannot get in speed %d: %v", s, err)
		}
		if s1 != s {
			t.Fatalf("In speed: %d != %d", s1, s)
		}
	}
	for _, s := range spds_fail {
		if err := ti.SetOSpeed(s); err != syscall.EINVAL {
			t.Fatalf("Set bad out speed %d", s)
		}
		if err := ti.SetISpeed(s); err != syscall.EINVAL {
			t.Fatalf("Set bad in speed %d", s)
		}
	}
}

var dev = os.Getenv("TEST_SERIAL_DEV")

func TestSetGet(t *testing.T) {
	if dev == "" {
		t.Skip("No TEST_SERIAL_DEV variable set.")
	}
	f, err := os.Open(dev)
	if err != nil {
		t.Fatal("Open:", err)
	}
	defer f.Close()

	var ti termios.Termios
	err = ti.GetFd(int(f.Fd()))
	if err != nil {
		t.Fatal("Cannot get termios:", err)
	}
	err = ti.SetOSpeed(19200)
	if err != nil {
		t.Fatal("Cannot set O speed to 19200:", err)
	}
	ti.SetISpeed(19200)
	if err != nil {
		t.Fatal("Cannot set I speed to 19200:", err)
	}
	err = ti.SetFd(int(f.Fd()), termios.TCSAFLUSH)
	if err != nil {
		t.Fatal("Cannot set termios:", err)
	}
	ti = termios.Termios{}
	err = ti.GetFd(int(f.Fd()))
	if err != nil {
		t.Fatal("Cannot get termios:", err)
	}
	spd, err := ti.GetOSpeed()
	if err != nil {
		t.Fatal("Cannot get O speed:", err)
	}
	if spd != 19200 {
		t.Fatalf("Bad O speed: %d != %d", spd, 19200)
	}
	spd, err = ti.GetISpeed()
	if err != nil {
		t.Fatal("Cannot get I speed:", err)
	}
	if spd != 19200 && spd != 0 {
		t.Fatalf("Bad I speed: %d != %d", spd, 19200)
	}
}

func TestMisc(t *testing.T) {
	if dev == "" {
		t.Skip("No TEST_SERIAL_DEV variable set.")
	}
	f, err := os.OpenFile(dev, os.O_RDWR, 0)
	if err != nil {
		t.Fatal("Open:", err)
	}
	defer f.Close()

	var ti termios.Termios
	err = ti.GetFd(int(f.Fd()))
	if err != nil {
		t.Fatal("Cannot get termios:", err)
	}
	/* Set low baudrate */
	baudrate := 1200
	ti.SetOSpeed(baudrate)
	ti.SetISpeed(baudrate)
	/* Disable flow control */
	ti.CFlag().Clr(termios.CRTSCTS)
	ti.IFlag().Clr(termios.IXON | termios.IXOFF | termios.IXANY)
	err = ti.SetFd(int(f.Fd()), termios.TCSANOW)
	if err != nil {
		t.Fatal("Cannot set termios:", err)
	}

	/* Try to test Drain */
	b := make([]byte, 600)
	start := time.Now()
	if _, err := f.Write(b); err != nil {
		t.Fatal("Cannot write:", err)
	}
	err = termios.Drain(int(f.Fd()))
	if err != nil {
		t.Fatal("Cannot drain:", err)
	}
	dur := time.Since(start)
	charTime := 10 * time.Second / time.Duration(baudrate)
	chars := int(dur / charTime)
	// Allow some fuzz for h/w queues and stuff.
	if chars < len(b)-16 || chars > len(b)+16 {
		t.Logf("Invalid tx time %v (%d chars):", dur, chars)
	}

	/* Try to test SendBreak */
	start = time.Now()
	termios.SendBreak(int(f.Fd()))
	dur = time.Since(start)
	// POSIX says SendBreak should last between 0.25 and 1 Sec.
	if dur < 200*time.Millisecond || dur > 1100*time.Millisecond {
		t.Log("Bad SendBreak duration:", dur)
	}

	// Just call Flush
	if err := termios.Flush(int(f.Fd()), termios.TCIFLUSH); err != nil {
		t.Fatal("Flush In failed:", err)
	}
	if err := termios.Flush(int(f.Fd()), termios.TCOFLUSH); err != nil {
		t.Fatal("Flush Out failed:", err)
	}
	if err := termios.Flush(int(f.Fd()), termios.TCIOFLUSH); err != nil {
		t.Fatal("Flush InOut failed:", err)
	}
	// This should normally fail (depends on system tcflush()
	// implementation)
	if err := termios.Flush(int(f.Fd()), 4242); err == nil {
		t.Logf("Flush 4242 should fail!")
	}
}
