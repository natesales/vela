package main

import "crypto/sha1"

func Checksum(data []byte) [20]byte {
	return sha1.Sum(data)
}
