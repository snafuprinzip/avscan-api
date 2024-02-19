package main

import (
	"log"
	"net/http"
	"net/netip"
	"strings"
)

// getHeaderIP retrieves the client's IP address from the request headers.
func getHeaderIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip, _, _ = strings.Cut(r.RemoteAddr, ":")
	}
	return ip
}

// isAccessGrantedByIP returns true if the client ip address is within the network ranges of the
// networks configured in the server.access section of the config file
func isAccessGrantedByIP(r *http.Request) (granted bool) {
	granted = false

	ip := getHeaderIP(r)
	ipaddr, err := netip.ParseAddr(ip)
	if err != nil {
		log.Printf("Unable to parse ip address %s: %s", ipaddr, err)
	}

	for _, allowed := range Config.Server.Access {
		network, err := netip.ParsePrefix(allowed)
		if err != nil {
			log.Printf("Unable to parse network prefix %s: %s", allowed, err)
		}
		if network.Contains(ipaddr) {
			granted = true
			break
		}
	}

	return
}
