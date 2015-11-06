// Copyright (c) 2015, Nick Patavalis (npat@efault.net).
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.txt file.

// +build linux freebsd netbsd openbsd darwin dragonfly solaris

// Implementation that uses the "termios" package for configuring the
// serial port, and the "poller" package for doing serial I/O. Should
// work with most Unix-like systems.

package serial

import (
	"strconv"
	"time"

	"github.com/npat-efault/poller"
	"github.com/npat-efault/serial/termios"
)

type port struct {
	fd          *poller.FD
	origTermios termios.Termios
	noReset     bool
}

func open(name string) (p *port, err error) {
	fd, err := poller.Open(name, poller.O_RW)
	if err != nil {
		return nil, newErr("open: " + err.Error())
	}

	if err := fd.Lock(); err != nil {
		return nil, ErrClosed
	}
	defer fd.Unlock()

	// Get attributes
	var tiosOrig termios.Termios
	err = tiosOrig.GetFd(fd.Sysfd())
	if err != nil {
		return nil, newErr("tcgetattr: " + err.Error())
	}
	// ?? Set HUPCL ??
	// tiosOrig.CFlag().Set(termios.HUPCL)
	// err = tiosOrig.SetFd(fd.Sysfd(), termios.TCSANOW)
	// if err != nil {
	// 	return nil, newErr("tcsetattr: " + err.Error())
	// }
	noReset := !tiosOrig.CFlag().Any(termios.HUPCL)

	// Set raw mode
	tios := tiosOrig
	tios.MakeRaw()
	err = tios.SetFd(fd.Sysfd(), termios.TCSANOW)
	if err != nil {
		return nil, newErr("tcsetattr: " + err.Error())
	}

	return &port{fd: fd, origTermios: tiosOrig, noReset: noReset}, nil
}

func (p *port) close() error {
	var errSetattr error

	if err := p.fd.Lock(); err != nil {
		return ErrClosed
	}
	defer p.fd.Unlock()

	if !p.noReset {
		err := p.origTermios.SetFd(p.fd.Sysfd(), termios.TCSANOW)
		if err != nil {
			errSetattr = newErr("tcsetattr: " + err.Error())
		}
	}
	err := p.fd.CloseUnlocked()
	if errSetattr != nil {
		err = errSetattr
	} else {
		if err != nil {
			err = newErr("close: " + err.Error())
		}
	}
	return err
}

func (p *port) getConf() (conf Conf, err error) {
	var tios termios.Termios
	var noReset bool

	if err = p.fd.Lock(); err != nil {
		return conf, ErrClosed
	}
	err = tios.GetFd(p.fd.Sysfd())
	noReset = p.noReset
	p.fd.Unlock()
	if err != nil {
		return conf, newErr("tcgetattr: " + err.Error())
	}

	// Baudrate
	conf.Baudrate, err = tios.GetOSpeed()
	if err != nil {
		return conf, newErr("getospeed: " + err.Error())
	}

	// Databits
	switch tios.CFlag().Msk(termios.CSIZE) {
	case termios.CS5:
		conf.Databits = 5
	case termios.CS6:
		conf.Databits = 6
	case termios.CS7:
		conf.Databits = 7
	case termios.CS8:
		conf.Databits = 8
	default:
		return conf, newErr("cannot decode databits")
	}

	// Stopbits
	if tios.CFlag().Any(termios.CSTOPB) {
		conf.Stopbits = 2
	} else {
		conf.Stopbits = 1
	}

	// Parity
	if !tios.CFlag().Any(termios.PARENB) {
		conf.Parity = ParityNone
	} else if tios.CFlag().Any(termios.CMSPAR) {
		if tios.CFlag().Any(termios.PARODD) {
			conf.Parity = ParityMark
		} else {
			conf.Parity = ParitySpace
		}
	} else {
		if tios.CFlag().Any(termios.PARODD) {
			conf.Parity = ParityOdd
		} else {
			conf.Parity = ParityEven
		}
	}

	// Flow
	rtscts := tios.CFlag().Any(termios.CRTSCTS)
	xoff := tios.IFlag().Any(termios.IXOFF)
	xon := tios.IFlag().Any(termios.IXON | termios.IXANY)

	if rtscts && !xoff && !xon {
		conf.Flow = FlowRTSCTS
	} else if !rtscts && xoff && xon {
		conf.Flow = FlowXONXOFF
	} else if !rtscts && !xoff && !xon {
		conf.Flow = FlowNone
	} else {
		conf.Flow = FlowOther
	}

	// NoReset
	conf.NoReset = noReset

	return conf, nil
}

func (p *port) confSome(conf Conf, flags ConfFlags) error {
	if err := p.fd.Lock(); err != nil {
		return ErrClosed
	}
	defer p.fd.Unlock()

	var tios termios.Termios
	err := tios.GetFd(p.fd.Sysfd())
	if err != nil {
		return newErr("tcgetattr: " + err.Error())
	}

	if flags&ConfBaudrate != 0 {
		err := tios.SetOSpeed(conf.Baudrate)
		if err != nil {
			return newErr("setospeed: " + err.Error())
		}
	}

	if flags&ConfDatabits != 0 {
		switch conf.Databits {
		case 5:
			tios.CFlag().Clr(termios.CSIZE).Set(termios.CS5)
		case 6:
			tios.CFlag().Clr(termios.CSIZE).Set(termios.CS6)
		case 7:
			tios.CFlag().Clr(termios.CSIZE).Set(termios.CS7)
		case 8:
			tios.CFlag().Clr(termios.CSIZE).Set(termios.CS8)
		default:
			return newErr("invalid databits value: " +
				strconv.Itoa(conf.Databits))
		}
	}

	if flags&ConfStopbits != 0 {
		switch conf.Stopbits {
		case 1:
			tios.CFlag().Clr(termios.CSTOPB)
		case 2:
			tios.CFlag().Set(termios.CSTOPB)
		default:
			return newErr("invalid stopbits value: " +
				strconv.Itoa(conf.Stopbits))
		}
	}

	if flags&ConfParity != 0 {
		switch conf.Parity {
		case ParityEven:
			tios.CFlag().Clr(termios.PARODD | termios.CMSPAR)
			tios.CFlag().Set(termios.PARENB)
		case ParityOdd:
			tios.CFlag().Clr(termios.CMSPAR)
			tios.CFlag().Set(termios.PARENB | termios.PARODD)
		case ParityMark:
			if termios.CMSPAR == 0 {
				return newErr("ParityMark not supported")
			}
			tios.CFlag().Set(termios.PARENB | termios.PARODD |
				termios.CMSPAR)
		case ParitySpace:
			if termios.CMSPAR == 0 {
				return newErr("ParitySpace not supported")
			}
			tios.CFlag().Clr(termios.PARODD)
			tios.CFlag().Set(termios.PARENB | termios.CMSPAR)
		case ParityNone:
			tios.CFlag().Clr(termios.PARENB | termios.PARODD |
				termios.CMSPAR)
		default:
			return newErr("invalid parity mode: " +
				conf.Parity.String())
		}
	}

	if flags&ConfFlow != 0 {
		switch conf.Flow {
		case FlowRTSCTS:
			if termios.CRTSCTS == 0 {
				return newErr("FlowRTSCTS not supported")
			}
			tios.CFlag().Set(termios.CRTSCTS)
			tios.IFlag().Clr(termios.IXON | termios.IXOFF |
				termios.IXANY)
		case FlowXONXOFF:
			tios.CFlag().Clr(termios.CRTSCTS)
			tios.IFlag().Set(termios.IXON | termios.IXOFF)
		case FlowNone:
			tios.CFlag().Clr(termios.CRTSCTS)
			tios.IFlag().Clr(termios.IXON | termios.IXOFF |
				termios.IXANY)
		default:
			return newErr("invalid flow-control mode: " +
				conf.Flow.String())
		}
	}

	if flags&ConfNoReset != 0 {
		p.noReset = conf.NoReset
		if p.noReset {
			tios.CFlag().Clr(termios.HUPCL)
		} else {
			tios.CFlag().Set(termios.HUPCL)
		}
	}

	err = tios.SetFd(p.fd.Sysfd(), termios.TCSANOW)
	if err != nil {
		return newErr("tcsetattr: " + err.Error())
	}

	return nil
}

func (p *port) read(b []byte) (n int, err error) {
	n, err = p.fd.Read(b)
	switch err {
	case poller.ErrTimeout:
		err = ErrTimeout
	case poller.ErrClosed:
		err = ErrClosed
	}
	return n, err
}

func (p *port) write(b []byte) (n int, err error) {
	n, err = p.fd.Write(b)
	switch err {
	case poller.ErrTimeout:
		err = ErrTimeout
	case poller.ErrClosed:
		err = ErrClosed
	}
	return n, err
}

func (p *port) setDeadline(t time.Time) error {
	err := p.fd.SetDeadline(t)
	if err == poller.ErrClosed {
		err = ErrClosed
	}
	return err
}

func (p *port) setReadDeadline(t time.Time) error {
	err := p.fd.SetReadDeadline(t)
	if err == poller.ErrClosed {
		err = ErrClosed
	}
	return err
}

func (p *port) setWriteDeadline(t time.Time) error {
	err := p.fd.SetReadDeadline(t)
	if err == poller.ErrClosed {
		err = ErrClosed
	}
	return err
}

func (p *port) flush(q flushSel) error {
	var qsel int
	switch q {
	case flushIn:
		qsel = termios.TCIFLUSH
	case flushOut:
		qsel = termios.TCOFLUSH
	case flushInOut:
		qsel = termios.TCIOFLUSH
	default:
		return newErr("invalid flush selector")
	}

	if err := p.fd.Lock(); err != nil {
		return ErrClosed
	}
	defer p.fd.Unlock()
	err := termios.Flush(p.fd.Sysfd(), qsel)
	if err != nil {
		return newErr("tcflush: " + err.Error())
	}
	return nil
}
