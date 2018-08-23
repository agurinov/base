package tools

import (
	"net"
	"strings"
)

func IsBot(ua string) bool {
	for _, word := range []string{
		"DuckDuckBot",
		"Slurp",
		"yahoo",
		"bingbot",
		"baiduspider",
		"yandex", "YandexBot",
		"yeti",
		"yodaobot",
		"gigabot",
		"ia_archiver",
		"facebookexternalhit", "Facebook", "facebot",
		"twitterbot",
		"LinkedInBot",
		"Googlebot",
		"msnbot",
		"flipboard",
		"MicroMessenger",
		"applebot",
	} {
		if strings.Contains(strings.ToLower(ua), strings.ToLower(word)) {
			return true
		}
	}

	return false
}

func GetRemoteAddr(some interface{}) net.Addr {
	if conn, ok := some.(net.Conn); ok {
		return conn.RemoteAddr()
	}

	return nil
}

func GetRemoteIP(addr net.Addr, headers ...string) net.IP {
	if ip := headersIP(headers...); ip != nil {
		return ip
	}

	return remoteIP(addr)
}

func isPrivateIP(ip net.IP) bool {
	for _, cidr := range []string{
		"127.0.0.0/8",    // IPv4 loopback
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
		"::1/128",        // IPv6 loopback
		"fe80::/10",      // IPv6 link-local
	} {
		if _, block, _ := net.ParseCIDR(cidr); block.Contains(ip) {
			return true
		}
	}

	return false
}

func headersIP(headers ...string) net.IP {
	for _, h := range headers {
		// LEFT TO RIGHT - current working variant
		for _, ipStr := range strings.Split(h, ",") {
			// header can contain spaces too, strip those out.
			ip := net.ParseIP(strings.TrimSpace(ipStr))
			if !ip.IsGlobalUnicast() || isPrivateIP(ip) {
				// bad address, go to next
				continue
			}
			return ip
		}

		// TODO: RIGHT TO LEFT - deal wi it (look at test)
		// TODO NOTE https://husobee.github.io/golang/ip-address/2015/12/17/remote-ip-go.html
		// TODO NOTE https://blog.gingerlime.com/2012/rails-ip-spoofing-vulnerabilities-and-protection/
		// addresses := strings.Split(h, ",")
		// // march from right to left until we get a public address
		// // that will be the address right before our proxy.
		// for i := len(addresses) - 1; i >= 0; i-- {
		// 	// header can contain spaces too, strip those out.
		// 	ip := net.ParseIP(strings.TrimSpace(addresses[i]))
		// 	if !ip.IsGlobalUnicast() || isPrivateIP(ip) {
		// 		// bad address, go to next
		// 		continue
		// 	}
		// 	return ip
		// }
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
