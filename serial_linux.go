// Copyright (c) 2015, Nick Patavalis (npat@efault.net).
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.txt file.

// +build linux

package serial

/*
#include <termios.h>
#include <unistd.h>
*/
import "C"
import "github.com/npat-efault/poller"

type port struct {
	fd          *poller.FD
	origTermios C.struct_termios
	noReset     bool
}

var stdSpeeds = speedTable{
	{0, C.B0},
	{50, C.B50},
	{75, C.B75},
	{110, C.B110},
	{134, C.B134},
	{150, C.B150},
	{200, C.B200},
	{300, C.B300},
	{600, C.B600},
	{1200, C.B1200},
	{1800, C.B1800},
	{2400, C.B2400},
	{4800, C.B4800},
	{9600, C.B9600},
	{19200, C.B19200},
	{38400, C.B38400},
	{57600, C.B57600},
	{115200, C.B115200},
	{230400, C.B230400},
	{460800, C.B460800},
	{500000, C.B500000},
	{576000, C.B576000},
	{921600, C.B921600},
	{1000000, C.B1000000},
	{1152000, C.B1152000},
	{2000000, C.B2000000},
	{2500000, C.B2500000},
	{3000000, C.B3000000},
	{3500000, C.B3500000},
	{4000000, C.B4000000},
}

func open(name string) (p *port, err error) {
	fd, err := poller.Open(name, poller.O_RW)
	if err != nil {
		return nil, &Error{ErrOpen, err}
	}

	if err := fd.Lock(); err != nil {
		return nil, &Error{ErrClosed, nil}
	}
	defer fd.Unlock()

	// Get attributes
	cfd := C.int(fd.Sysfd())
	var tiosOrig C.struct_termios
	r, err := C.tcgetattr(cfd, &tiosOrig)
	if r < 0 {
		return nil, &Error{ErrGetConf, err}
	}

	// Set raw mode, CLOCAL and HUPCL
	tios := tiosOrig
	C.cfmakeraw(&tios)
	tios.c_cflag |= C.CLOCAL | C.HUPCL
	r, err = C.tcsetattr(cfd, C.TCSANOW, &tios)
	if r < 0 {
		return nil, &Error{ErrConf, err}
	}

	return &port{fd: fd, origTermios: tiosOrig}, nil
}

func (p *port) close() error {
	var errSetattr error
	if !p.noReset {
		if err := p.fd.Lock(); err != nil {
			return &Error{ErrClosed, nil}
		}
		r, err := C.tcsetattr(C.int(p.fd.Sysfd()),
			C.TCSANOW, &p.origTermios)
		p.fd.Unlock()
		if r < 0 {
			errSetattr = &Error{ErrReset, err}
		} else {
			errSetattr = nil
		}
	}
	err := p.fd.Close()
	if errSetattr != nil {
		return errSetattr
	} else {
		return err
	}
}

func (p *port) getConf() (conf Conf, err error) {
	var tios C.struct_termios

	if err := p.fd.Lock(); err != nil {
		return conf, &Error{ErrClosed, nil}
	}
	r, err := C.tcgetattr(C.int(p.fd.Sysfd()), &tios)
	p.fd.Unlock()
	if r < 0 {
		return conf, &Error{ErrGetConf, err}
	}

	// Baudrate
	spdCode := C.cfgetospeed(&tios)
	var ok bool
	conf.Baudrate, ok = stdSpeeds.Speed(uint32(spdCode))
	if !ok {
		return conf, &Error{ErrGetConf, nil}
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
		return conf, &Error{ErrGetConf, nil}
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
	conf.NoReset = p.noReset

	return conf, nil
}

func (p *port) doConf(conf Conf, flags int) error {
	var tios C.struct_termios

	if err := p.fd.Lock(); err != nil {
		return &Error{ErrClosed, nil}
	}
	defer p.fd.Unlock()

	cfd := C.int(p.fd.Sysfd())
	r, err := C.tcgetattr(cfd, &tios)
	if r < 0 {
		return &Error{ErrGetConf, err}
	}

	if flags&dcBaudrate != 0 {
		spd, ok := stdSpeeds.Code(conf.Baudrate)
		if !ok {
			return &Error{ErrInvalid, nil}
		}
		C.cfsetospeed(&tios, C.speed_t(spd))
	}

	if flags&dcDatabits != 0 {
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
			return &Error{ErrInvalid, nil}
		}
	}

	if flags&dcStopbits != 0 {
		switch conf.Stopbits {
		case 1:
			tios.c_cflag &^= C.CSTOPB
		case 2:
			tios.c_cflag |= C.CSTOPB
		default:
			return &Error{ErrInvalid, nil}
		}
	}

	if flags&dcParity != 0 {
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
			return &Error{ErrInvalid, nil}
		}
	}

	if flags&dcFlow != 0 {
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
			return &Error{ErrInvalid, nil}
		}
	}

	if flags&dcNoReset != 0 {
		p.noReset = conf.NoReset
	}

	r, err = C.tcsetattr(cfd, C.TCSANOW, &tios)
	if r < 0 {
		return &Error{ErrConf, err}
	}

	return nil
}
