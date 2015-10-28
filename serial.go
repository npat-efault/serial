package serial

// Port is a serial port
type Port struct {
	Name string
	*port
}

// ParityMode encodes the supported bit-parity modes
type ParityMode int

const (
	ParityNone  ParityMode = iota // No parity bit
	ParityEven                    // Even bit-parity
	ParityOdd                     // Odd bit-parity
	ParityMark                    // Parity bit to logical 1 (mark)
	ParitySpace                   // Parity bit to logical 0 (space)
)

// FlowMode encodes the supported flow-control modes
type FlowMode int

const (
	FlowNone    FlowMode = iota // No flow control
	FlowRTSCTS                  // Hardware flow control
	FlowXONXOFF                 // Software flow control
	FlowOther                   // Unknown mode
)

// Conf is used to pass the serial port's configuration parameters to
// and from methods of this package.
type Conf struct {
	Baudrate int        // in Bits Per Second
	Databits int        // 5, 6, 7, or 8
	Stopbits int        // 1 or 2
	Parity   ParityMode // see ParityXXX constants
	Flow     FlowMode   // see FlowXXX constants
	NoReset  bool       // don't reset and don't hangup on close
}

// Open opens the named serial port
func Open(name string) (port *Port, err error) {
	p, err := open(name)
	if err != nil {
		return nil, err
	}
	return &Port{Name: name, port: p}, nil
}

// Close closes the port.
func (p *Port) Close() error {
	return p.port.close()
}

// GetConf returns the serial port's configuration parameters as a
// Conf structure.
func (p *Port) GetConf() (conf Conf, err error) {
	return p.port.getConf()
}

// doConf flags, controlling which parameters to configure
const (
	dcBaudrate = 1 << iota
	dcDatabits
	dcStopbits
	dcParity
	dcFlow
	dcNoReset
	dcAll = dcBaudrate | dcDatabits | dcStopbits |
		dcParity | dcFlow | dcNoReset
)

func (p *Port) doConf(conf Conf, flags int) error {
	return p.port.doConf(conf, flags)
}

// Conf configures the serial port using the parameters in the Conf
// structure
func (p *Port) Conf(conf Conf) error {
	return p.doConf(conf, dcAll)
}

func (p *Port) SetBaudrate(b int) error {
	conf := Conf{Baudrate: b}
	return p.doConf(conf, dcBaudrate)
}

func (p *Port) SetDatabits(d int) error {
	conf := Conf{Databits: d}
	return p.doConf(conf, dcDatabits)
}

func (p *Port) SetStopbits(s int) error {
	conf := Conf{Stopbits: s}
	return p.doConf(conf, dcStopbits)
}

func (p *Port) SetParity(r ParityMode) error {
	conf := Conf{Parity: r}
	return p.doConf(conf, dcParity)
}

func (p *Port) SetFlow(f FlowMode) error {
	conf := Conf{Flow: f}
	return p.doConf(conf, dcFlow)
}

func (p *Port) SetNoReset(nr bool) {
	p.port.noReset = nr
}

/*
func (p *Port) Read(b []byte) (n int, err error) {
	return p.port.read(b)
}

func (p *Port) Write(b []byte) (n int, err error) {
	return p.port.write(b)
}
*/

// speedTable is used to map numeric tty speeds (baudrates) to the
// respective code (Bxxx) values.
type speedTable []struct {
	speed int
	code  uint32
}

func (t speedTable) Code(speed int) (code uint32, ok bool) {
	for _, s := range t {
		if s.speed == speed {
			return s.code, true
		}
	}
	return 0, false
}

func (t speedTable) Speed(code uint32) (speed int, ok bool) {
	for _, s := range t {
		if s.code == code {
			return s.speed, true
		}
	}
	return 0, false
}
