// Copyright (c) 2015, Nick Patavalis (npat@efault.net).
// All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.

// !! Tests are necessarily superficial !!

package termios_test

import (
	"os"
	"syscall"
	"testing"

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
	spds_ok := []int{9600, 19200, 38400, 115200}
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
		t.Fatal("Cannot set speed to 19200:", err)
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
		t.Fatal("Cannot get speed:", err)
	}
	if spd != 19200 {
		t.Fatalf("Bad speed: %d != %d", spd, 19200)
	}
}
