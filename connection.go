package main

import "net"

type Circuit struct {
	remote net.IP
	mtu    uint16
	vid    uint8
}

func NewCircuit(remote net.IP, vid uint8) Circuit {
	return Circuit{
		remote: remote,
		mtu:    9000,
		vid:    vid,
	}
}
