package ip

import (
	"errors"
	"net"
	"net/http"
	"strings"
)

func GetIpFromRequest(r *http.Request) (string, error) {
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}

	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip, nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}
	return "", errors.New("No IP information found")
}

// GetRemoteAddr parses the given request, resolves the X-Forwarded-For header
// and returns the resolved remote address.
func GetRemoteAddr(r *http.Request, numProxies int) string {
	if xffh := r.Header.Get("X-Forwarded-For"); xffh != "" {
		if sip, sport, err := net.SplitHostPort(r.RemoteAddr); err == nil && sip != "" {
			if xip := Parse(xffh, numProxies); xip != "" {
				return net.JoinHostPort(xip, sport)
			}
		}
	}
	return r.RemoteAddr
}

// Parse parses the value of the X-Forwarded-For Header and returns the IP address.
func Parse(xffString string, numProxies int) string {
	ipList := strings.Split(xffString, ",")
	if numProxies <= 0 {
		numProxies = len(ipList)
	}
	for idx := len(ipList) - 1; idx >= 0; idx-- {
		ipStr := strings.TrimSpace(ipList[idx])
		if ip := net.ParseIP(ipStr); ip != nil && IsPublicIP(ip) {
			if numProxies <= 0 || idx == 0 {
				return ipStr
			}
			numProxies--
		}
	}
	return ""
}
