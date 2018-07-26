package tools

import (
	"net"
)

func GetRemoteIP() net.IP {
	return net.ParseIP("185.86.151.11")
}
