package main

import (
	"log"
	"net"
)

func InitListener() *net.UDPConn {
	// Inbound Listener
	listenerAddr, err := net.ResolveUDPAddr("udp", ":48621")
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

	return listenerConn
}
