package main

const HEADER_SIZE = 2 // Header size in bytes

type Header struct {
	VC  byte // VELA Control Code
	VID byte // VELA Circuit ID
}

func (header Header) Parse() []byte {
	return []byte{header.VC, header.VID}
}
