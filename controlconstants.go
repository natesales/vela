package main

// VC (VELA Control) Codes
const (
	VC_NOP  byte = 0 // No Operation
	VC_IREQ byte = 1 // Session Initialization Request
	VC_IACK byte = 2 // Session Initialization Acknowledgement
	VC_ICON byte = 3 // Session Initialization Confirmation
	VC_CREQ byte = 4 // Session Closure Request
)
