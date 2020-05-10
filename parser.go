package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

func main() {
	data, err := ioutil.ReadFile("vela.yml")
	if err != nil {
		log.Fatalln("File reading error", err)
	}

	circuits := make(map[string]map[string]string)

	err = yaml.Unmarshal(data, &circuits)
	if err != nil {
		log.Fatalln("Could not parse config file: ", err)
	}

	for iface, params := range circuits {
		if !strings.HasPrefix(iface, "vela") {
			log.Fatalln("Interface " + iface + " not identified as VELA circuit.")
		} else {
			var id uint8
			var remote net.IP
			var port uint16
			var network net.IP
			var mask uint8
			var ipv6 bool

			_id, err := strconv.Atoi(params["id"])
			if err != nil {
				log.Fatalln("Circuit ID must be an integer:", err)
			} else {
				if _id > 255 {
					log.Fatalln("Circuit ID must be less than 256")
				} else {
					id = uint8(_id)
				}
			}

			_remote := strings.Split(params["remote"], ":")
			if len(_remote) == 2 { // If no port specified
				_port, err := strconv.Atoi(_remote[1])
				if err != nil {
					log.Fatalln("Error parsing remote UDP port:", err)
				} else {
					if _port > 65535 {
						log.Fatalln("Remote UDP port must not be larger than 65535")
					} else {
						port = uint16(_port)
					}
				}
			} else if len(_remote) < 1 || len(_remote) > 2 {
				log.Fatalln("Remote must be an IPv4 or IPv6 address with or without a port.")
			} else {
				port = 1337 // Default port
			}

			remote = net.ParseIP(_remote[0])
			if remote == nil {
				log.Fatalln("Remote must be an IPv4 or IPv6 address with or without a port.")
			}

			_network := strings.Split(params["network"], "/")
			if len(_network) != 2 {
				log.Fatalln("Network address must be an IP with subnet mask")
			} else {
				network = net.ParseIP(_network[0])
				if network == nil {
					log.Fatalln("Remote must be an IPv4 or IPv6 address")
				}

				if network.To4() != nil {
					ipv6 = false
				} else {
					ipv6 = true
				}
			}

			_mask, err := strconv.Atoi(_network[1])
			if err != nil {
				log.Fatalln("Can't parse netmask:", err)
			} else {
				if !ipv6 && _mask >= 32 {
					log.Fatalln("Invalid IPv4 netmask:", _mask)
				} else if ipv6 && _mask >= 127 {
					log.Fatalln("Invalid IPv6 netmask:", _mask)
				} else {
					mask = uint8(_mask)
				}
			}

			if ipv6 {
				fmt.Print("IPv6")
			} else {
				fmt.Print("IPv4")
			}

			fmt.Printf(" tunnel to %s:%v with id %v carrying %v/%v\n", remote, port, id, network, mask)
		}
	}
}
