package tools

import (
	"net"
	"strings"
)

func GetRemoteAddr(some interface{}) net.Addr {
	if conn, ok := some.(net.Conn); ok {
		return conn.RemoteAddr()
	}

	return nil
}

// https://husobee.github.io/golang/ip-address/2015/12/17/remote-ip-go.html
func GetRemoteIP(addr net.Addr, headers ...string) net.IP {
	if ip := headersIP(headers...); ip != nil {
		return ip
	}

	return remoteIP(addr)
}

// func headersIP(headers ...string) net.IP {
// 	for _, h := range headers {
// 		addresses := strings.Split(h, ",")
// 		// march from right to left until we get a public address
// 		// that will be the address right before our proxy.
// 		for i := len(addresses) - 1; i >= 0; i-- {
// 			log.Debug("IP FOR CHECK: ", strings.TrimSpace(addresses[i]))
// 			// header can contain spaces too, strip those out.
// 			ipStr := strings.TrimSpace(addresses[i])
// 			if ip := net.ParseIP(ipStr); ip.IsGlobalUnicast() {
// 				return ip
// 			}
// 			// bad address, go to next
// 		}
// 	}
//
// 	return nil
// }

func headersIP(headers ...string) net.IP {
	for _, h := range headers {
		for _, ipStr := range strings.Split(h, ",") {
			// header can contain spaces too, strip those out.
			if ip := net.ParseIP(strings.TrimSpace(ipStr)); ip.IsGlobalUnicast() {
				return ip
			}
			// bad address, go to next
		}
	}

	return nil
}

func remoteIP(addr net.Addr) net.IP {
	switch typed := addr.(type) {
	case *net.UDPAddr:
		return typed.IP
	case *net.TCPAddr:
		return typed.IP
	default:
		return nil
	}
}
