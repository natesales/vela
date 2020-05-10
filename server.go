package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
)

// Temporary configuration constants
const (
	host         string = "localhost"
	port         uint16 = 1337
	usev6        bool   = false
	sanityChecks bool   = false // TODO: Enable pre-release
)

var preprocessor func([]byte) []byte = nil
var postprocessor func([]byte) []byte = nil

// Protocol-specific constants
const (
	mtu = 1500
)

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func main() {
	var addr net.IP

	_addr := net.ParseIP(host)
	if _addr == nil { // If not valid address
		addresses, err := net.LookupIP(host)
		if err != nil { // If DNS resolution failed
			fmt.Println("Unknown host")
			os.Exit(1)
		} else { // If DNS resolution succeeded
			for _, address := range addresses {
				if usev6 {
					if address.To4() == nil { // If IPv6 forced
						addr = address
						break
					}
				} else { // If IPv4 forced
					if address.To4() != nil { // If IPv4
						addr = address
						break
					}
				}
			}

			if addr == nil { // If an address can't be found
				fmt.Println("No address for host.")
				os.Exit(1)
			}
		}
	} else {
		addr = _addr
	}

	if sanityChecks && !addr.IsGlobalUnicast() {
		fmt.Printf("Address is not globally routable, are you sure you want to continue? [Y/n] ")

		reader := bufio.NewReader(os.Stdin)
		if confirm, _ := reader.ReadString('\n'); !(confirm == "Y\n" || confirm == "\n" || confirm == "y\n") {
			fmt.Println("Exiting")
			os.Exit(0)
		}
	}

	var networkProtocol string
	if !usev6 {
		networkProtocol = "udp4"
	} else {
		networkProtocol = "udp6"
	}

	s, err := net.ResolveUDPAddr(networkProtocol, ":"+strconv.Itoa(int(port))) // Verify UDP address resolution
	if err != nil {
		fmt.Println(err)
		return
	}

	connection, err := net.ListenUDP(networkProtocol, s)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		_ = connection.Close() // Ignore errors on closing since we're going to exit anyway
	}()

	buffer := make([]byte, mtu)

	fmt.Println("Ready on " + host + ":" + strconv.Itoa(int(port)))

	for {
		n, addr, err := connection.ReadFromUDP(buffer)
		incoming := buffer[0:n]

		vc := incoming[0]
		vid := incoming[1]

		fmt.Print(addr, " sent ")

		switch vc {
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
			fmt.Println("Unknown VC" + strconv.Itoa(int(vc)))
		}

		_, err = connection.WriteToUDP([]byte{'0'}, addr)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
