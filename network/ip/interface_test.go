package ip

import (
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/duke-git/lancet/v2/netutil"
)

// TestGetInternalIp Get internal ip information.
// func GetInternalIp() string
func TestGetInternalIp(t *testing.T) {
	internalIp := netutil.GetInternalIp()
	ip := net.ParseIP(internalIp)
	fmt.Println(ip) // 127.0.0.1
}

// TestGetIps Get all ipv4 list.
// func GetIps() string
func TestGetIps(t *testing.T) {
	ips := netutil.GetIps()
	fmt.Println(ips) // [192.168.2.161 192.168.255.10]
}

// TestGetMacAddrs Get all ipv4 list.
// func GetMacAddrs() []string {
func TestGetMacAddrs(t *testing.T) {
	macAddrs := netutil.GetMacAddrs()
	fmt.Println(macAddrs)
}

// TestGetPublicIpInfo Get public ip information.
// func GetPublicIpInfo() (*PublicIpInfo, error)
//
//	type PublicIpInfo struct {
//	    Status      string  `json:"status"`
//	    Country     string  `json:"country"`
//	    CountryCode string  `json:"countryCode"`
//	    Region      string  `json:"region"`
//	    RegionName  string  `json:"regionName"`
//	    City        string  `json:"city"`
//	    Lat         float64 `json:"lat"`
//	    Lon         float64 `json:"lon"`
//	    Isp         string  `json:"isp"`
//	    Org         string  `json:"org"`
//	    As          string  `json:"as"`
//	    Ip          string  `json:"query"`
//	}
func TestGetPublicIpInfo(t *testing.T) {
	publicIpInfo, err := netutil.GetPublicIpInfo()
	if err != nil {
		fmt.Println(err)
	}

	// &{success China CN GD Guangdong Shenzhen 22.5559 114.0577 Chinanet Chinanet GD AS4134 CHINANET-BACKBONE 11.11.11.11}
	fmt.Println(publicIpInfo)
}

// TestGetRequestPublicIp Get http request public ip.
// func GetRequestPublicIp(req *http.Request) string
func TestGetRequestPublicIp(t *testing.T) {
	ip := "36.112.24.10"

	request := http.Request{
		Method: "GET",
		Header: http.Header{
			"X-Forwarded-For": {ip},
		},
	}
	publicIp := netutil.GetRequestPublicIp(&request)

	fmt.Println(publicIp)
}
