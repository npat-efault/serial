// Copyright (c) 2015, Nick Patavalis (npat@efault.net).
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.txt file.

// +build linux freebsd netbsd openbsd darwin dragonfly solaris

// POSIX implementation that uses CGo ang LIBC's termios functions.

package serial

/*

#include <termios.h>
#include <unistd.h>

#ifndef CRTSCTS
#define CRTSCTS 0
#endif

#ifndef CMSPAR
#define CMSPAR 0
#endif

*/
import "C"
import (
	"strconv"
	"syscall"
	"time"

	"github.com/npat-efault/poller"
)

type port struct {
	fd          *poller.FD
	origTermios C.struct_termios
	noReset     bool
}

func tcSetAttr(fd, act C.int, tios *C.struct_termios) (n C.int, err error) {
	for {
		n, err := C.tcsetattr(fd, act, tios)
		if n < 0 && err == syscall.EINTR {
			continue
		}
		return n, err
	}
}

func tcGetAttr(fd C.int, tios *C.struct_termios) (n C.int, err error) {
	for {
		n, err := C.tcgetattr(fd, tios)
		if n < 0 && err == syscall.EINTR {
			continue
		}
		return n, err
	}
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
	cfd := C.int(fd.Sysfd())
	var tiosOrig C.struct_termios
	r, err := tcGetAttr(cfd, &tiosOrig)
	if r < 0 {
		return nil, newErr("tcgetattr: " + err.Error())
	}

	// Set raw mode, CLOCAL and HUPCL
	tios := tiosOrig
	cfMakeRaw(&tios)
	tios.c_cflag |= C.CLOCAL | C.HUPCL
	r, err = tcSetAttr(cfd, C.TCSANOW, &tios)
	if r < 0 {
		return nil, newErr("tcsetattr: " + err.Error())
	}

	return &port{fd: fd, origTermios: tiosOrig}, nil
}

func (p *port) close() error {
	var errSetattr error

	if err := p.fd.Lock(); err != nil {
		return ErrClosed
	}
	defer p.fd.Unlock()

	if !p.noReset {
		r, err := tcSetAttr(C.int(p.fd.Sysfd()),
			C.TCSANOW, &p.origTermios)
		if r < 0 {
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
	var tios C.struct_termios
	var noReset bool

	if err := p.fd.Lock(); err != nil {
		return conf, ErrClosed
	}
	r, err := tcGetAttr(C.int(p.fd.Sysfd()), &tios)
	noReset = p.noReset
	p.fd.Unlock()
	if r < 0 {
		return conf, newErr("tcgetattr: " + err.Error())
	}

	// Baudrate
	spdCode := C.cfgetospeed(&tios)
	var ok bool
	conf.Baudrate, ok = stdSpeeds.Speed(uint32(spdCode))
	if !ok {
		return conf, newErr("cannot decode baudrate")
	}

	// Databits
	switch tios.c_cflag & C.CSIZE {
	case C.CS5:
		conf.Databits = 5
	case C.CS6:
		conf.Databits = 6
	case C.CS7:
		conf.Databits = 7
	case C.CS8:
		conf.Databits = 8
	default:
		return conf, newErr("cannot decode databits")
	}

	// Stopbits
	if tios.c_cflag&C.CSTOPB == 0 {
		conf.Stopbits = 1
	} else {
		conf.Stopbits = 2
	}

	// Parity
	flg := tios.c_cflag
	if flg&C.PARENB == 0 {
		conf.Parity = ParityNone
	} else if flg&C.CMSPAR != 0 {
		if flg&C.PARODD != 0 {
			conf.Parity = ParityMark
		} else {
			conf.Parity = ParitySpace
		}
	} else {
		if flg&C.PARODD != 0 {
			conf.Parity = ParityOdd
		} else {
			conf.Parity = ParityEven
		}
	}

	// Flow
	rtscts := tios.c_cflag&C.CRTSCTS != 0
	xoff := tios.c_iflag&C.IXOFF != 0
	xon := tios.c_iflag&(C.IXON|C.IXANY) != 0

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
	var tios C.struct_termios

	if err := p.fd.Lock(); err != nil {
		return ErrClosed
	}
	defer p.fd.Unlock()

	if flags&ConfNoReset != 0 {
		p.noReset = conf.NoReset
	}
	if flags & ^ConfNoReset == 0 {
		return nil
	}

	cfd := C.int(p.fd.Sysfd())
	r, err := tcGetAttr(cfd, &tios)
	if r < 0 {
		return newErr("tcgetattr: " + err.Error())
	}

	if flags&ConfBaudrate != 0 {
		spd, ok := stdSpeeds.Code(conf.Baudrate)
		if !ok {
			return newErr("invalid baudrate: " +
				strconv.Itoa(conf.Baudrate))
		}
		C.cfsetospeed(&tios, C.speed_t(spd))
	}

	if flags&ConfDatabits != 0 {
		switch conf.Databits {
		case 5:
			tios.c_cflag &^= C.CSIZE
			tios.c_cflag |= C.CS5
		case 6:
			tios.c_cflag &^= C.CSIZE
			tios.c_cflag |= C.CS6
		case 7:
			tios.c_cflag &^= C.CSIZE
			tios.c_cflag |= C.CS7
		case 8:
			tios.c_cflag &^= C.CSIZE
			tios.c_cflag |= C.CS8
		default:
			return newErr("invalid databits value: " +
				strconv.Itoa(conf.Databits))
		}
	}

	if flags&ConfStopbits != 0 {
		switch conf.Stopbits {
		case 1:
			tios.c_cflag &^= C.CSTOPB
		case 2:
			tios.c_cflag |= C.CSTOPB
		default:
			return newErr("invalid stopbits value: " +
				strconv.Itoa(conf.Stopbits))
		}
	}

	if flags&ConfParity != 0 {
		switch conf.Parity {
		case ParityEven:
			tios.c_cflag &^= C.PARODD | C.CMSPAR
			tios.c_cflag |= C.PARENB
		case ParityOdd:
			tios.c_cflag &^= C.CMSPAR
			tios.c_cflag |= C.PARENB | C.PARODD
		case ParityMark:
			tios.c_cflag |= C.PARENB | C.PARODD | C.CMSPAR
		case ParitySpace:
			tios.c_cflag &^= C.PARODD
			tios.c_cflag |= C.PARENB | C.CMSPAR
		case ParityNone:
			tios.c_cflag &^= C.PARENB | C.PARODD | C.CMSPAR
		default:
			return newErr("invalid parity mode: " +
				conf.Parity.String())
		}
	}

	if flags&ConfFlow != 0 {
		switch conf.Flow {
		case FlowRTSCTS:
			tios.c_cflag |= C.CRTSCTS
			tios.c_iflag &^= C.IXON | C.IXOFF | C.IXANY
		case FlowXONXOFF:
			tios.c_cflag &^= C.CRTSCTS
			tios.c_iflag |= C.IXON | C.IXOFF
		case FlowNone:
			tios.c_cflag &^= C.CRTSCTS
			tios.c_iflag &^= C.IXON | C.IXOFF | C.IXANY
		default:
			return newErr("invalid flow-control mode: " +
				conf.Flow.String())
		}
	}

	r, err = tcSetAttr(cfd, C.TCSANOW, &tios)
	if r < 0 {
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

func tcFlush(fd C.int, qsel C.int) (n C.int, err error) {
	for {
		n, err := C.tcflush(fd, qsel)
		if n < 0 && err == syscall.EINTR {
			continue
		}
		return n, err
	}
}

func (p *port) flush(q flushSel) error {
	var qsel C.int
	switch q {
	case flushIn:
		qsel = C.TCIFLUSH
	case flushOut:
		qsel = C.TCOFLUSH
	case flushInOut:
		qsel = C.TCIOFLUSH
	default:
		return newErr("invalid flush selector")
	}

	if err := p.fd.Lock(); err != nil {
		return ErrClosed
	}
	defer p.fd.Unlock()
	r, err := tcFlush(C.int(p.fd.Sysfd()), qsel)
	if r < 0 {
		err = newErr("tcflush: " + err.Error())
		return err
	}
	return nil
}
