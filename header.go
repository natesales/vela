package main

const HEADER_SIZE = 2 // Header size in bytes TODO: use reflect to calculate this automatically from the struct below

type Header struct {
	VC  byte // VELA Control Code
	VID byte // VELA Circuit ID
}

func (header Header) Parse() []byte {
	return []byte{header.VC, header.VID}
}
