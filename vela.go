package main

import (
	"fmt"
	"github.com/songgao/water"
	"log"
	"net"
)

var (
	remote string = "10.0.0.27:48621"
)

var (
	BUFFERSIZE = 9000                    // TODO: make this the interface mtu + overhead
	MTU        = BUFFERSIZE - 8 - 1 - 20 // 8-byte UDP header, 1-byte VELA Circuit ID, and 20 byte IP header
)

func main() {
	listenerConn := InitListener()

	circuits := ParseConfig()

	var links map[*net.UDPAddr]*water.Interface

	for _, circuit := range circuits {
		SetInterface(circuit.Iface, circuit.Network, MTU)

		// Resolve remote address
		remote, err := net.ResolveUDPAddr("udp", circuit.Remote)
		if nil != err {
			log.Fatalln("Unable to resolve remote addr:", err)
		}

		go func() {
			// Transmit in primary Goroutine
			packet := make([]byte, BUFFERSIZE)
			for {
				packetLength, err := circuit.Iface.Read(packet)
				if err != nil {
					break
				}

				pkt := append([]byte{2}, packet[:packetLength]...) // Prepend Circuit ID

				_, err = listenerConn.WriteToUDP(pkt, remote)
				if err != nil {
					log.Println("Error sending data over UDP connection:", err)
				}
			}
		}()
	} // End circuit enumeration

	// Receive Goroutine
	// TODO: automatic buffer adjustment from link.Attrs().MTU

	buffer := make([]byte, BUFFERSIZE)

	for {
		n, addr, err := listenerConn.ReadFromUDP(buffer)
		vid := buffer[0]

		if err != nil || n == 0 {
			fmt.Println("Error: ", err)
		}

		if remoteUDPConn, ok := links[addr]; ok { // If address exists in map
			go func() {
				_, _ = remoteUDPConn.Write(buffer[1:n]) // Only write the payload to the interface TODO: If error in transmit, send a message requesting the previous packet
			}()
		} // If address isn't in the map, discard packet.

		fmt.Println(addr, "vid", vid)
	}
}
