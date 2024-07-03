package main

import "net"

func mustParseTrustedSubnetCIDR(cidrString string) *net.IPNet {
	if cidrString == "" {
		return nil
	}

	_, ipnet, err := net.ParseCIDR(cidrString)
	if err != nil {
		panic(err)
	}
	return ipnet
}
