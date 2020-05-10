package v1

import (
	"fmt"
	"strconv"
)

// VC (VELA Control) Codes
const (
	VC_NOP  byte = 0 // No Operation
	VC_IREQ byte = 1 // Session Initialization Request
	VC_IACK byte = 2 // Session Initialization Acknowledgement
	VC_ICON byte = 3 // Session Initialization Confirmation
	VC_CREQ byte = 4 // Session Closure Request
)

func ParseVC(header Header) {
	switch header.VC {
	case VC_NOP:
		fmt.Println("VC-NOP")
	case VC_IREQ:
		fmt.Println("VC-IREQ")
	case VC_IACK:
		fmt.Println("VC-IACK")
	case VC_ICON:
		fmt.Println("VC-ICON")
	case VC_CREQ:
		fmt.Println("VC-CREQ")
	default: // If VC is unknown
		fmt.Println("Unknown VC" + strconv.Itoa(int(header.VC)))
	}
}
