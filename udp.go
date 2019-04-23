package udpfacade

import (
	"syscall"
	"net"
	"github.com/google/gopacket"
	//"github.com/google/gopacket/layers"
	"github.com/eric-lindau/gopacket/layers"
	"encoding/binary"
	"unsafe"
	"runtime"
)

// TODO Determine if pool of sockets might be effective

// A lot of credit to Graham King for this fantastic article:
// https://www.darkcoding.net/software/raw-sockets-in-go-link-layer/

// NOTE: Requires sudo for raw socket
func WriteTo(b []byte, src *net.IPAddr, dst *net.UDPAddr) error {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	if err != nil {
		return err
	}

	if err := syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_HDRINCL, 1); err != nil {
		return err
	}

	to := &syscall.SockaddrInet4{
		Addr: [4]byte{dst.IP[0], dst.IP[1], dst.IP[2], dst.IP[3]},
		Port: 0,
	}

	p := gopacket.NewSerializeBuffer()
	if err := craftPacket(b, &p, &src.IP, dst); err != nil {
		return err
	}

	if err := syscall.Sendto(fd, p.Bytes(), 0, to); err != nil {
		return err
	}

	syscall.Close(fd) // TODO Ensure conn closes quickly

	return nil
}

func craftPacket(b []byte, p *gopacket.SerializeBuffer, src *net.IP, dst *net.UDPAddr) error {
	opts := gopacket.SerializeOptions{
		FixLengths: true,
	}

	ipv4 := layers.IPv4{
		Version:    4,
		IHL:        5,
		FragOffset: 16,
		TTL:        64,
		Protocol:   layers.IPProtocolUDP,
		SrcIP:      *src,
		DstIP:      (*dst).IP,
	}

	udp := layers.UDP{
		DstPort: layers.UDPPort(dst.Port),
	}

	if err := gopacket.SerializeLayers(*p, opts, &ipv4, &udp, gopacket.Payload(b)); err != nil {
		panic(err)
	}

	switch runtime.GOOS { // NOTE: freebsd < version 11 not supported
	case "darwin", "dragonfly", "openbsd":
		ipv4.RawSocketByteOrder(*p, nativeEndian)
	}

	return nil
}

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
