package main

import (
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/songgao/water"
	"gopkg.in/yaml.v2"
)

type CircuitConfig struct {
	Iface   *water.Interface // VELA interface
	Remote  string           // Remote IP and port
	Vid     uint8            // Circuit ID
	Network string           // IP address with mask
	IPv6    bool             // Is remote IPv6?
}

func ParseConfig() []CircuitConfig {
	data, err := ioutil.ReadFile("vela.yml")
	if err != nil {
		log.Fatalln("File reading error", err)
	}

	_circuits := make(map[string]map[string]string)

	err = yaml.Unmarshal(data, &_circuits)
	if err != nil {
		log.Fatalln("Could not parse config file: ", err)
	}

	circuits := make([]CircuitConfig, 0)

	for _iface, params := range _circuits {
		if !strings.HasPrefix(_iface, "vela") {
			log.Fatalln("Interface " + _iface + " not identified as VELA circuit.")
		}
		var id uint8
		var remote net.IP
		var port string // Type is string to avoid further parsing for compatibility with the UDP resolver
		var network net.IP
		var mask string // Type is string to avoid further parsing for compatibility with the UDP resolver
		var ipv6 bool
		var iface *water.Interface

		iface, err := water.New(water.Config{
			DeviceType:             water.TUN,
			PlatformSpecificParams: water.PlatformSpecificParams{Name: _iface},
		})
		if err != nil {
			log.Fatalln("Error creating interface")
		}

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
					port = strconv.Itoa(_port)
				}
			}
		} else if len(_remote) < 1 || len(_remote) > 2 {
			log.Fatalln("Remote must be an IPv4 or IPv6 address with or without a port.")
		} else {
			port = strconv.Itoa(PORT) // Default port, type is string to avoid further parsing for compatibility with the UDP resolver
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
				mask = strconv.Itoa(_mask)
			}
		}

		circuits = append(circuits, CircuitConfig{
			Iface:   iface,
			Remote:  remote.String() + ":" + port,
			Vid:     id,
			Network: network.String() + "/" + mask,
			IPv6:    ipv6,
		})
	}

	return circuits
}
