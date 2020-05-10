package main

import (
	"fmt"
	"github.com/songgao/water"
	"github.com/vishvananda/netlink"
	"log"
	"net"
	"os"
	"strconv"
)

// Usage: vela [local network] [remote]

var (
	BUFFERSIZE = 9000
	MTU        = BUFFERSIZE - 8 - 1 // 8-byte UDP header and 1-byte VELA Circuit ID
	localIP    = os.Args[1]   // IP with mask
	remoteIP   = os.Args[2]
	port       = 4321
)

func main() {
	var iface *water.Interface
	var err error                 // Explicit declare error for multiple assignment
	for i := 0; i <= 32768; i++ { // (const int max_netdevices = 8*PAGE_SIZE in __dev_alloc_name) TODO: Extract to constant
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

	// Get interface objects
	link, _ := netlink.LinkByName(iface.Name())
	addr, _ := netlink.ParseAddr(localIP)

	// Configure the link
	_ = netlink.LinkSetMTU(link, MTU)
	_ = netlink.AddrAdd(link, addr)
	_ = netlink.LinkSetUp(link)

	// Resolve remote address
	remote, err := net.ResolveUDPAddr("udp", remoteIP+":"+strconv.Itoa(port))
	if nil != err {
		log.Fatalln("Unable to resolve remote addr:", err)
	}

	// Inbound Listener
	listenerAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port)) // TODO: Explicit local listen address?
	if nil != err {
		log.Fatalln("Unable to resolve listener UDP socket:", err)
	}
	listenerConn, err := net.ListenUDP("udp", listenerAddr)
	if nil != err {
		log.Fatalln("Unable to listen UDP socket:", err)
	}
	defer func() {
		_ = listenerConn.Close()
	}()

	// Receive Goroutine
	go func() {
		buffer := make([]byte, BUFFERSIZE)

		for {
			n, addr, err := listenerConn.ReadFromUDP(buffer)
			vid := buffer[0]

			if err != nil || n == 0 {
				fmt.Println("Error: ", err)
			}

			_, _ = iface.Write(buffer[1:n]) // Only write the payload to the interface TODO: If error in transmit, send a message requesting the previous packet

			fmt.Println(addr, "vid", vid)
		}
	}()

	// Transmit in primary Goroutine
	packet := make([]byte, BUFFERSIZE)
	for {
		packetLength, err := iface.Read(packet)
		if err != nil {
			break
		}

		pkt := append([]byte{2}, packet[:packetLength]...) // Prepend Circuit ID

		_, err = listenerConn.WriteToUDP(pkt, remote)
		if err != nil {
			log.Println("Error sending data over UDP connection:", err)
		}
	}
}