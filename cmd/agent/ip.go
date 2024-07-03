package main

import (
	"log"
	"net"
)

func mustGetLocalIPAddr(serverEndpoint string) *net.IP {
	conn, err := net.Dial("udp", serverEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return &localAddr.IP
}
