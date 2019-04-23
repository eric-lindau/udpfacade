package udpfacade

import (
	"testing"
	"net"
)

const (
	dstPort = 1025
)

// NOTE: Requires sudo
func TestWriteTo(t *testing.T) {
	spoof := &net.IPAddr{
		IP: net.IPv4(1, 2, 3, 4),
	}
	self := &net.IPAddr{
		IP: net.IPv4(127, 0, 0, 1),
	}

	err := WriteTo([]byte("a"),
		&net.UDPAddr{
			IP: spoof.IP, Port: 1234, Zone: "",
		},
		&net.UDPAddr{
			IP: self.IP, Port: dstPort, Zone: "",
		})

	if err != nil {
		t.Error(err)
	}
}
