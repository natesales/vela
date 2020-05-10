package main

const HEADER_SIZE = 2 // Header size in bytes

type Header struct {
	vc  byte // VELA Control Code
	vid byte // VELA Circuit ID
}

func (header Header) Parse() []byte {
	return []byte{header.vc, header.vid}
}
