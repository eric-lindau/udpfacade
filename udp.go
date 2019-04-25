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
)

// IMPLEMENT LATER: sync Pool for efficiency

type UDPConn struct {
	sock int
	to   *syscall.SockaddrInet4
	Src  *net.UDPAddr
	Dst  *net.UDPAddr
}

// Sets up connection using the src properties for outgoing UDP/IP headers.
// NOTE: Requires sudo for raw socket
func DialUDPFrom(src *net.UDPAddr, dst *net.UDPAddr) (*UDPConn, error) {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	if err != nil {
		return nil, err
	}

	if err := syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_HDRINCL, 1); err != nil {
		return nil, err
	}

	to := &syscall.SockaddrInet4{
		Addr: [4]byte{dst.IP[0], dst.IP[1], dst.IP[2], dst.IP[3]},
		Port: 0,
	}

	return &UDPConn{sock: fd, to: to, Src: src, Dst: dst}, nil
}

func (c *UDPConn) Write(b []byte) error {
	p := gopacket.NewSerializeBuffer()
	if err := craftPacket(b, &p, c.Src, c.Dst); err != nil {
		return err
	}

	if err := syscall.Sendto(c.sock, p.Bytes(), 0, c.to); err != nil {
		return err
	}

	return nil
}

func (c *UDPConn) Close() error {
	if err := syscall.Close(c.sock); err != nil {
		return err
	}

	return nil
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
