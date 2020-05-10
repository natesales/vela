package main

import (
	"fmt"
	"net"
)

func main() {
	CONNECT := "localhost:1337"

	s, err := net.ResolveUDPAddr("udp4", CONNECT)
	c, err := net.DialUDP("udp4", nil, s)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("The UDP server is %s\n", c.RemoteAddr().String())
	defer c.Close() // TODO: Error handling

	_, err = c.Write([]byte{0, 0})

	//
	//buffer := make([]byte, 1024)
	//n, _, err := c.ReadFromUDP(buffer)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Printf("Reply: %s\n", string(buffer[0:n]))
}
