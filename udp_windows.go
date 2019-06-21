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

func (c *UDPConn) Read(b []byte) (int, error) {
	return -1, errors.New("cannot read from transparent udp connection")
}

func (c *UDPConn) Write(b []byte) (int, error) {
	return -1, nil
}

func (c *UDPConn) Close() error {
    return nil
}

func (c *UDPConn) LocalAddr() net.Addr {
    return nil
}

func (c *UDPConn) RemoteAddr() net.Addr {
    return nil
}

func (c *UDPConn) SetDeadline(t time.Time) error {
    return nil
}

func (c *UDPConn) SetReadDeadline(t time.Time) error {
    return nil
}

func (c *UDPConn) SetWriteDeadline(t time.Time) error {
    return nil
}

func craftPacket(b []byte, p *gopacket.SerializeBuffer, src *net.UDPAddr, dst *net.UDPAddr) error {
	return nil
}
