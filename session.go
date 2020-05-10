package main

import "net"

var sessions map[string]map[byte]net.IP

func Set(source net.IP, vid byte, remote net.IP) bool {
	if sessions[string(source)][vid] == nil { // If session doesn't already exist
		sessions[string(source)][vid] = remote // Register the session
		return true
	} else { // If the session doesn't exist
		return false // Don't do anything and report
	}
}

func Get(source net.IP, vid byte) net.IP {
	return sessions[string(source)][vid]
}

func Delete(source net.IP, vid byte) bool {
	delete(sessions[string(source)], vid) // TODO: Test this. Otherwise set to nil
}
