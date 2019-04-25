package udpfacade

import (
	"syscall"
	"net"
	"github.com/google/gopacket"
	//"github.com/google/gopacket/layers" --- Currently bugged for some raw IPv4 sockets
	"github.com/eric-lindau/gopacket/layers"
	"encoding/binary"
	"unsafe"
	"runtime"
	"fmt"
	"errors"
	"time"
)

// IMPLEMENT LATER: sync Pool for efficiency

// Transparent UDP connection
type UDPConn struct {
	conn *net.IPConn
	Src  *net.UDPAddr
	Dst  *net.UDPAddr
}

// Sets up connection using the src properties for outgoing UDP/IP headers.
// NOTE: Requires sudo for raw socket
func DialUDPFrom(src *net.UDPAddr, dst *net.UDPAddr) (*UDPConn, error) {
	conn, err := net.DialIP(fmt.Sprintf("ip:%d", syscall.IPPROTO_RAW), nil, &net.IPAddr{
		IP: src.IP, Zone: "",
	})
	if err != nil {
		return nil, err
	}

	f, err := conn.File()
	fd := f.Fd()
	if err := syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IP, syscall.IP_HDRINCL, 1); err != nil {
		return nil, err
	}

	return &UDPConn{conn: conn, Src: src, Dst: dst}, nil
}

func (c *UDPConn) Read(b []byte) (int, error) {
	return -1, errors.New("cannot read from transparent udp connection")
}

func (c *UDPConn) Write(b []byte) (int, error) {
	p := gopacket.NewSerializeBuffer()
	if err := craftPacket(b, &p, c.Src, c.Dst); err != nil {
		return 0, err
	}

	return c.conn.Write(p.Bytes())
}

func (c *UDPConn) Close() error {
	return c.conn.Close()
}

func (c *UDPConn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *UDPConn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *UDPConn) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

func (c *UDPConn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *UDPConn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

func craftPacket(b []byte, p *gopacket.SerializeBuffer, src *net.UDPAddr, dst *net.UDPAddr) error {
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	ipv4 := layers.IPv4{
		Version:  4,
		IHL:      5,
		TTL:      64,
		Protocol: layers.IPProtocolUDP,
		SrcIP:    (*src).IP,
		DstIP:    (*dst).IP,
	}

	udp := layers.UDP{
		SrcPort: layers.UDPPort(src.Port),
		DstPort: layers.UDPPort(dst.Port),
	}
	udp.SetNetworkLayerForChecksum(&ipv4)

	if err := gopacket.SerializeLayers(*p, opts, &ipv4, &udp, gopacket.Payload(b)); err != nil {
		return err
	}

	switch runtime.GOOS { // NOTE: freebsd < version 11 not supported
	case "darwin", "dragonfly", "openbsd":
		ipv4.RawSocketByteOrder(*p, nativeEndian)
	}

	return nil
}

// Native endianness detection for syscalls on some platforms
// https://github.com/tensorflow/tensorflow/blob/master/tensorflow/go/tensor.go#L488-L505
var nativeEndian binary.ByteOrder

func init() {
	buf := [2]byte{}
	*(*uint16)(unsafe.Pointer(&buf[0])) = uint16(0xABCD)

	switch buf {
	case [2]byte{0xCD, 0xAB}:
		nativeEndian = binary.LittleEndian
	case [2]byte{0xAB, 0xCD}:
		nativeEndian = binary.BigEndian
	default:
		panic("Could not determine native endianness.")
	}
}
