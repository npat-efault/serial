package serial

import (
	"os"
	"testing"
)

var dev = os.Getenv("TEST_SERIAL_DEV")

func TestBaudrate(t *testing.T) {
	if dev == "" {
		t.Skip("No TEST_SERIAL_DEV variable set.")
	}
	p, err := Open(dev)
	if err != nil {
		t.Fatal("Open:", err)
	}
	c0, err := p.GetConf()
	if err != nil {
		t.Fatal("GetConf:", err)
	}

	for _, spd := range stdSpeeds {
		s := spd.speed
		c := Conf{Baudrate: s}
		err := p.ConfSome(c, ConfBaudrate)
		if err != nil {
			if s <= 38400 {
				// All devices must go up to 38400
				t.Fatalf("ConfSome, Baudrate %v: %v", s, err)
			} else {
				t.Logf("ConfSome, Baudrate %v: %v (OK?)",
					s, err)
				continue
			}
		}
		c, err = p.GetConf()
		if err != nil {
			t.Fatalf("GetConf, Baudrate %v: %v", s, err)
		}
		if c.Baudrate != s {
			if s <= 38400 {
				// All devices must go up to 38400
				t.Fatalf("Baudrate: %d != %d", c.Baudrate, s)
			} else {
				t.Logf("Baudrate: %d != %d (OK?)",
					c.Baudrate, s)
			}
		}
		c.Baudrate = c0.Baudrate
		if c != c0 {
			t.Fatalf("%v != %v", c, c0)
		}
	}

	err = p.Close()
	if err != nil {
		t.Fatal("Close:", err)
	}
}

func TestDatsbits(t *testing.T) {
	if dev == "" {
		t.Skip("No TEST_SERIAL_DEV variable set.")
	}
	p, err := Open(dev)
	if err != nil {
		t.Fatal("Open:", err)
	}
	c0, err := p.GetConf()
	if err != nil {
		t.Fatal("GetConf:", err)
	}

	databits := []int{5, 6, 7, 8}

	for _, y := range databits {
		c := Conf{Databits: y}
		err := p.ConfSome(c, ConfDatabits)
		if err != nil {
			if y == 5 || y == 6 {
				// Some devices do not support 5 or 6 db
				t.Logf("ConfSome, Databits %v: %v (OK?)",
					y, err)
				continue
			} else {
				t.Fatalf("ConfSome, Databits %v: %v", y, err)
			}
		}
		c, err = p.GetConf()
		if err != nil {
			t.Fatalf("GetConf, Databits %v: %v", y, err)
		}
		if c.Databits != y {
			if y == 5 || y == 6 {
				// Some devices do not support 5 or 6 db
				t.Logf("Databits: %v != %v (OK?)",
					c.Databits, y)
			} else {
				t.Fatalf("Databits: %v != %v", c.Databits, y)
			}
		}
		c.Databits = c0.Databits
		if c != c0 {
			t.Fatalf("%v != %v", c, c0)
		}
	}

	err = p.Close()
	if err != nil {
		t.Fatal("Close:", err)
	}
}

func TestParity(t *testing.T) {
	if dev == "" {
		t.Skip("No TEST_SERIAL_DEV variable set.")
	}
	p, err := Open(dev)
	if err != nil {
		t.Fatal("Open:", err)
	}
	c0, err := p.GetConf()
	if err != nil {
		t.Fatal("GetConf:", err)
	}

	parities := []ParityMode{ParityNone, ParityEven, ParityOdd,
		ParityMark, ParitySpace}

	for _, y := range parities {
		c := Conf{Parity: y}
		err := p.ConfSome(c, ConfParity)
		if err != nil {
			if y == ParityMark || y == ParitySpace {
				t.Logf("ConfSome, Parity %v: %v", y, err)
				continue
			} else {
				t.Fatalf("ConfSome, Parity %v: %v", y, err)
			}
		}
		c, err = p.GetConf()
		if err != nil {
			t.Fatalf("GetConf, Parity %v: %v", y, err)
		}
		if c.Parity != y {
			if y == ParityMark || y == ParitySpace {
				// Some systems do not support Mark and Space
				t.Logf("Parity: %v != %v (OK?)",
					c.Parity, y)
				continue
			} else {
				t.Fatalf("Parity: %v != %v", c.Parity, y)
			}
		}
		c.Parity = c0.Parity
		if c != c0 {
			t.Fatalf("%v != %v", c, c0)
		}
	}

	err = p.Close()
	if err != nil {
		t.Fatal("Close:", err)
	}
}

func TestStopbits(t *testing.T) {
	if dev == "" {
		t.Skip("No TEST_SERIAL_DEV variable set.")
	}
	p, err := Open(dev)
	if err != nil {
		t.Fatal("Open:", err)
	}
	c0, err := p.GetConf()
	if err != nil {
		t.Fatal("GetConf:", err)
	}

	stopbits := []int{1, 2}

	for _, y := range stopbits {
		c := Conf{Stopbits: y}
		err := p.ConfSome(c, ConfStopbits)
		if err != nil {
			t.Fatalf("ConfSome, Stopbits %v: %v", y, err)
		}
		c, err = p.GetConf()
		if err != nil {
			t.Fatalf("GetConf, Stopbits %v: %v", y, err)
		}
		if c.Stopbits != y {
			t.Fatalf("Stopbits: %v != %v", c.Stopbits, y)
		}
		c.Stopbits = c0.Stopbits
		if c != c0 {
			t.Fatalf("%v != %v", c, c0)
		}
	}

	err = p.Close()
	if err != nil {
		t.Fatal("Close:", err)
	}
}

func TestFlow(t *testing.T) {
	if dev == "" {
		t.Skip("No TEST_SERIAL_DEV variable set.")
	}
	p, err := Open(dev)
	if err != nil {
		t.Fatal("Open:", err)
	}
	c0, err := p.GetConf()
	if err != nil {
		t.Fatal("GetConf:", err)
	}

	flows := []FlowMode{FlowNone, FlowRTSCTS, FlowXONXOFF}

	for _, y := range flows {
		c := Conf{Flow: y}
		err := p.ConfSome(c, ConfFlow)
		if err != nil {
			t.Fatalf("ConfSome, Flow %v: %v", y, err)
		}
		c, err = p.GetConf()
		if err != nil {
			t.Fatalf("GetConf, Flow %v: %v", y, err)
		}
		if c.Flow != y {
			t.Fatalf("Flow: %v != %v", c.Flow, y)
		}
		c.Flow = c0.Flow
		if c != c0 {
			t.Fatalf("%v != %v", c, c0)
		}
	}

	err = p.Close()
	if err != nil {
		t.Fatal("Close:", err)
	}
}

func TestNoReset(t *testing.T) {
	if dev == "" {
		t.Skip("No TEST_SERIAL_DEV variable set.")
	}
	p, err := Open(dev)
	if err != nil {
		t.Fatal("Open:", err)
	}
	c0, err := p.GetConf()
	if err != nil {
		t.Fatal("GetConf:", err)
	}

	s := 9600
	if c0.Baudrate == 9600 {
		s = 19200
	}

	c := Conf{Baudrate: s, NoReset: true}
	err = p.ConfSome(c, ConfBaudrate|ConfNoReset)
	if err != nil {
		t.Fatal("ConfSome:", err)
	}

	err = p.Close()
	if err != nil {
		t.Fatal("Close:", err)
	}
	p, err = Open(dev)
	if err != nil {
		t.Fatal("Open:", err)
	}
	c, err = p.GetConf()
	if err != nil {
		t.Fatal("GetConf:", err)
	}
	if c.Baudrate != s {
		t.Fatalf("Baudrate %v != %v", c.Baudrate, s)
	}

	err = p.Close()
	if err != nil {
		t.Fatal("Close:", err)
	}
}
