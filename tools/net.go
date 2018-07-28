package tools

import (
	"net"
	"strings"
)

func GetRemoteIP(xForwardedFor string, remoteAddr net.Addr) net.IP {
	// firstly, check forwarded for from proxy load balancers
	ipStr := strings.SplitN(xForwardedFor, ",", 2)[0]
	if ip := net.ParseIP(ipStr); ip != nil {
		return ip
	}

	// TODO second part get remote addr of connection
	// return net.ParseIP("185.86.151.11")

	return nil
}
