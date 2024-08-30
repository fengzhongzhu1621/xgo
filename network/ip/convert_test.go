package ip

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInet4ToInt64(t *testing.T) {
	testCases := []struct {
		ip     net.IP
		expect int64
	}{
		{net.ParseIP("192.168.0.1"), 3232235521},
		{net.ParseIP("10.0.0.1"), 167772161},
	}

	for _, tc := range testCases {
		actual := Inet4ToInt64(tc.ip)
		assert.Equal(t, tc.expect, actual)
	}
}

func TestIpv4ToInt64(t *testing.T) {
	testCases := []struct {
		ip     string
		expect int64
	}{
		{"192.168.0.1", 3232235521},
		{"10.0.0.1", 167772161},
		{"256.256.256.256", 0}, // out of range IP address
	}

	for _, tc := range testCases {
		actual := Ipv4ToInt64(tc.ip)
		assert.Equal(t, tc.expect, actual)
	}
}

func TestInetToHexStr(t *testing.T) {
	cases := []struct {
		ip     net.IP
		expect string
	}{
		{net.ParseIP("2001:db8::1"), "20010db8000000000000000000000001"},
		{net.ParseIP("192.168.0.1"), "c0a80001"},
		{net.ParseIP("10.0.0.1"), "0a000001"},
	}

	for _, tc := range cases {
		actual := InetToHexStr(tc.ip)
		assert.Equal(t, tc.expect, actual)
	}
}

func TestCheckIpInCidrList(t *testing.T) {
	testCases := []struct {
		srcIp    string
		cidrList []string
		expect   bool
	}{
		{
			srcIp:    "192.168.1.1",
			cidrList: []string{"192.168.1.0/24", "10.10.0.0/16"},
			expect:   true,
		},
		{
			srcIp:    "10.0.0.1",
			cidrList: []string{"192.168.1.0/24", "10.10.0.0/16"},
			expect:   false,
		},
	}
	for _, tc := range testCases {
		actual := CheckIpInCidrList(tc.srcIp, tc.cidrList)
		assert.Equal(t, tc.expect, actual, "CheckIpInCidrList(%v, %v)", tc.srcIp, tc.cidrList)
	}
}
