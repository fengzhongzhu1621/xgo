package validator

import (
	"errors"
	"log"
	"net"
	"regexp"
)

var (
	privateIPBlocks []*net.IPNet
	privateCIDRs    = []string{
		"127.0.0.0/8",    // IPv4 loopback
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
		"169.254.0.0/16", // RFC3927 link-local
		"::1/128",        // IPv6 loopback
		"fe80::/10",      // IPv6 link-local
		"fc00::/7",       // IPv6 unique local addr
	}
)

func init() {
	for _, cidr := range privateCIDRs {
		switch _, block, err := net.ParseCIDR(cidr); {
		case err != nil:
			log.Fatalf("invalid cidr %q: %v", cidr, err)
		default:
			privateIPBlocks = append(privateIPBlocks, block)
		}
	}
}

func isPrivateIP(ip net.IP) bool {
	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

func IsPublicIP(ip net.IP) bool {
	return !ip.IsLoopback() && !ip.IsLinkLocalUnicast() && !ip.IsLinkLocalMulticast() && !isPrivateIP(ip)
}

func IsIPv6Addr(ip string) (bool, error) {
	if net.ParseIP(ip) == nil {
		return false, errors.New("Address parsing failed")
	}
	for i := 0; i < len(ip); i++ {
		switch ip[i] {
		case '.':
			return false, nil
		case ':':
			return true, nil
		}
	}
	return false, errors.New("Unable to determine address type")
}

func IsIp(ip string) bool {
	reg := regexp.MustCompile(`^(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|[1-9])\.(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)\.(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)\.(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)$`)
	return reg.MatchString(ip)
}
