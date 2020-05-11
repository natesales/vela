package main

import (
	"github.com/songgao/water"
	"github.com/vishvananda/netlink"
	"log"
)

func SetInterface(iface *water.Interface, address string, mtu int) {
	// Get interface objects
	link, _ := netlink.LinkByName(iface.Name())
	addr, _ := netlink.ParseAddr(address)

	// Configure the link
	err := netlink.LinkSetMTU(link, mtu)
	if err != nil {
		log.Fatalln("Error setting link attribute: ", err)
	}
	err = netlink.AddrAdd(link, addr)
	if err != nil {
		log.Fatalln("Error setting link attribute: ", err)
	}
	err = netlink.LinkSetUp(link)
	if err != nil {
		log.Fatalln("Error setting link attribute: ", err)
	}
}
