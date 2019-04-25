# udpfacade
Small library that might be useful for transparent UDP forwarding or UDP/IP spoofing, among other things.

Currently supports UDP over IPv4.

This library differs from most others than implement this functionality in that it is able to exploit the security/reliability of Go's built-in IP Dialing and Writing instead of writing frames directly to the wire or directly using raw IP sockets.

It took a long time to trace through the source to find a way to allow this, so I hope some find this useful!

## Resources
* Graham King's fantastic [article(s)](https://www.darkcoding.net/software/raw-sockets-in-go-link-layer/) about low-level networking in Go.
* Jan Newmarch's set of [pages](https://ipfs.io/ipfs/QmfYeDhGH9bZzihBUDEQbCbTc5k5FZKURMUoUvfmc27BwL/index.html) about networking in Go.
* [gopacket](https://godoc.org/github.com/google/gopacket)
