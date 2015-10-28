package serial

import "testing"

func TestOpen(t *testing.T) {
	p, err := Open("/dev/ttyUSB0")
	if err != nil {
		t.Fatal("Open:", err)
	}
	c, err := p.GetConf()
	if err != nil {
		t.Fatal("GetConf:", err)
	}
	t.Logf("Conf: %+v", c)

	err = p.Close()
	if err != nil {
		t.Fatal("Close:", err)
	}
}
