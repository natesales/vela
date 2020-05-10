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

	_, err = c.Write([]byte{0, 1}) //TODO: NEXT FIX: For some reason, if vid is 1 then it works, but if vid is 0 then the server thinks it's an invalid header.

	//
	//buffer := make([]byte, 1024)
	//n, _, err := c.ReadFromUDP(buffer)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Printf("Reply: %s\n", string(buffer[0:n]))
}
