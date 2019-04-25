package udpfacade

import (
	"testing"
	"net"
)

const (
	dstPort = 4321
)

// NOTE: Requires sudo. Wireshark recommended to see results (no server setup required).
func TestWriteTo(t *testing.T) {
	spoof := &net.UDPAddr{
		IP: net.IPv4(1, 2, 3, 4),
		Port: 1234,
	}
	self := &net.UDPAddr{
		IP: net.IPv4(127, 0, 0, 1),
		Port: dstPort,
	}

	conn, err := DialUDPFrom(spoof, self)
	if err != nil {
		t.Error(err)
	}

	if err := conn.Write([]byte("aa")); err != nil {
		t.Error(err)
	}
	if err := conn.Write([]byte("bb")); err != nil {
		t.Error(err)
	}

	if err := conn.Close(); err != nil {
		t.Error(err)
	}
}
