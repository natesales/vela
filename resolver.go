package main

import (
	"fmt"
	"net"
	"os"
)

func Resolve(host string, forcev6 bool) net.IP {
	addr := net.ParseIP(host)
	if addr == nil { // If not valid address
		addresses, err := net.LookupIP(host)
		if err != nil { // If DNS resolution failed
			fmt.Println("Unknown host")
			os.Exit(1)
		} else { // If DNS resolution succeeded
			for _, address := range addresses {
				if forcev6 {
					if address.To4() == nil { // If IPv6 forced
						addr = address
						break
					}
				} else {                      // If IPv4 forced
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
		return addr
	}

	return nil
}
