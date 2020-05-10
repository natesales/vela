package main

import (
	"fmt"
	"net"
)

func Write(connection net.UDPConn, source net.IP, vid byte, header Header, payload []byte) bool {
	if Get(source, vid) { // If session is established, send the data
		_, err := connection.WriteToUDP(append(header.Parse(), payload...), &net.UDPAddr{IP: source, Port: PORT})
		if err != nil { // If there's an error in sending
			fmt.Println(err)
		}
		return true
	} else { // If the session isn't established, don't send anything
		return false
	}
}

// [vc, vid, ...]
func Read(incoming []byte) (bool, Header, []byte) { // Returns header and payload
	if len(incoming) >= HEADER_SIZE { // If incoming is less than a valid header, then we can't do anything.
		fmt.Println("Header good, returning.")
		return true, Header{incoming[0], incoming[1]}, incoming[2:]
	} else {
		fmt.Println("Header too small.")
		return false, Header{}, nil // If incoming is invalid data, return nil values.
	}
}
