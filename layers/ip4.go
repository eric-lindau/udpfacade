package layers

import (
	"encoding/binary"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// For raw IP sockets on darwin, dragonfly, netbsd, and freebsd (before ver. 11), you must call
// this method with host byte order before writing b.Bytes() to the socket.
//
// On such systems, some packet fields written to raw IP sockets are expected in host byte order.
func RawSocketByteOrder(ip *layers.IPv4, b gopacket.SerializeBuffer, o binary.ByteOrder) {
	bytes := b.Bytes()[2:8]
	o.PutUint16(bytes[0:], ip.Length)
	o.PutUint16(bytes[4:], ip.FragOffset&0x1FFF|uint16(ip.Flags)<<13)
}

