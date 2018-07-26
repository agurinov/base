package tools

import (
	"net"
)

func GetRemoteIP() net.IP {
	// httpRequest.Header.Get("X-Forwarded-For") | httpRequest.Header.Get("X-Forwarded-For"))
	return net.ParseIP("185.86.151.11")
}
