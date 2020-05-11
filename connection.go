package main

import (
	"github.com/songgao/water"
	"log"
	"net"
	"strconv"
)

type Circuit struct {
	Remote net.IP
	Vid    uint8
	Iface  *water.Interface
}

func NewCircuit(remote net.IP, vid uint8) Circuit {
	var iface *water.Interface
	for i := 0; i <= 32768; i++ { // Max interfaces
		// Create tunnel interface
		_iface, err := water.New(water.Config{
			DeviceType:             water.TUN,
			PlatformSpecificParams: water.PlatformSpecificParams{Name: "vela" + strconv.Itoa(i)},
		})

		if err == nil { // If interface allocation succeeded, break the loop
			iface = _iface
			break
		}

		if err.Error() != "device or resource busy" {
			log.Fatalln("Unable to allocate TUN interface:", err)
		}
	}

	if iface == nil {
		log.Fatalln("Interface allocation failed.")
	} else {
		log.Println("Interface allocated:", iface.Name())
	}

	return Circuit{
		Remote: remote,
		Vid:    vid,
		Iface:  iface,
	}
}
