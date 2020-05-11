package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/songgao/water"
	"github.com/vishvananda/netlink"
)

// Usage vela [remote] [local] [network] [verbose]

var (
	remote  = os.Args[1] + ":48621"
	local   = os.Args[2] + ":48621"
	network = os.Args[3]
	mtu     = 9000
	verbose = true
)

type Circuit struct {
	remote  string
	local   string
	network string
	name    string
	mtu     uint16
}

const (
	Port          = 48621 // Default VELA UDP port
	MaxInterfaces = 100   // Maximum number of VELA interfaces on a system
)

var adjustedMTU = mtu - 8 - 1 - 20 // 8-byte UDP header, 1-byte VELA Circuit ID, and 20 byte IP header

func setMTU(link netlink.Link, mtu int) {
	err := netlink.LinkSetMTU(link, mtu)
	if err != nil {
		log.Fatalln("Error setting link attribute: ", err)
	}
}

func addAddr(link netlink.Link, addr *netlink.Addr) {
	err := netlink.AddrAdd(link, addr)
	if err != nil {
		log.Fatalln("Error setting link attribute: ", err)
	}
}

func linkUp(link netlink.Link) {
	err := netlink.LinkSetUp(link)
	if err != nil {
		log.Fatalln("Error setting link attribute: ", err)
	}
}

// Resolve an IPv{4/6} address or hostname and prefer IPv6 connections
func Resolve(host string) net.IP {
	addr := net.ParseIP(host)
	if addr == nil { // If not valid address
		addresses, err := net.LookupIP(host)
		if err != nil { // If DNS resolution failed
			log.Fatalln("Address resolution failed:", err)
		} else { // If DNS resolution succeeded
			if len(addresses) < 1 {
				log.Fatalln("No address for host", host)
			}

			for _, address := range addresses { // If we can find an IPv6 address, use that.
				if address.To4() == nil { // If address is IPv6
					addr = address
					break
				}
			}

			addr = addresses[0] // If the host is IPv4 only, then use that.
		}
	}

	return addr
}

func main() {

	if os.Args[4] == "true" {
		verbose = true
	} else {
		verbose = false
	}

	var iface *water.Interface
	for i := 0; i <= MaxInterfaces; i++ {
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

		// If it failed because the selected index is already in use, try the next sequential index
	}

	if iface == nil {
		log.Fatalln("Interface allocation failed.")
	}

	fmt.Println("Interface allocated:", iface.Name())

	// Get interface objects
	link, _ := netlink.LinkByName(iface.Name())
	addr, _ := netlink.ParseAddr(network)

	// Configure the link
	setMTU(link, adjustedMTU)
	addAddr(link, addr)
	linkUp(link)

	// Resolve remote addr
	remoteAddr, err := net.ResolveUDPAddr("udp", remote)
	if nil != err {
		log.Fatalln("Unable to resolve remote addr:", err)
	}
	// listen to local socket
	listenerAddr, err := net.ResolveUDPAddr("udp", local)
	if nil != err {
		log.Fatalln("Unable to get UDP socket:", err)
	}
	listenerConn, err := net.ListenUDP("udp", listenerAddr)
	if nil != err {
		log.Fatalln("Unable to listen on UDP socket:", err)
	}
	defer listenerConn.Close()

	// RX stream goroutine
	go func() {
		packet := make([]byte, mtu)
		for {
			packetLength, addr, _ := listenerConn.ReadFromUDP(packet) // Read full packet from UDP stream
			vid := packet[0]
			_, _ = iface.Write(packet[1:packetLength]) // Send to tunnel interface

			if verbose {
				log.Println("Received", packetLength-1, "bytes from", addr, "with vid", vid) // Subtract 1 for the actual packet length (accounting for the vid)
			}
		}
	}()

	// TX stream
	packet := make([]byte, mtu) // Packet buffer
	for {
		packetLength, err := iface.Read(packet) // Read from interface
		if err != nil {
			break
		}
		_, _ = listenerConn.WriteToUDP(append([]byte{1}, packet[:packetLength]...), remoteAddr) // Send to remote

		if verbose {
			log.Println("Sent", packetLength, "bytes to", remoteAddr.String())
		}
	}
}
