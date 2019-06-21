// +build windows

package udpfacade

import (
	"net"
	"errors"
)

// Transparent UDP connection
type UDPConn struct {
	conn *net.IPConn
	Src  *net.UDPAddr
	Dst  *net.UDPAddr
}

// Mock implementation to prevent Windows build errors - this should be implemented later
func DialUDPFrom(src *net.UDPAddr, dst *net.UDPAddr) (*UDPConn, error) {
	return nil, errors.New("transparent udp connection not implemented for windows")
}
