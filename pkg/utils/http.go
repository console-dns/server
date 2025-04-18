package utils

import (
	"net"
	"net/http"
	"net/netip"
)

func ParseAddress(r *http.Request) net.IP {
	ip := net.ParseIP(r.RemoteAddr)
	if ip != nil {
		return ip
	}
	addr, _ := netip.ParseAddrPort(r.RemoteAddr)
	return net.ParseIP(addr.Addr().String())
}
