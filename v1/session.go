package v1

import "net"

var sessions map[string]map[byte]bool

func Set(source net.IP, vid byte) {
	sessions[string(source)][vid] = true
}

func Get(source net.IP, vid byte) bool {
	return sessions[string(source)][vid]
}

func Delete(source net.IP, vid byte) {
	delete(sessions[string(source)], vid)
}
