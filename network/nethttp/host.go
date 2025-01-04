package nethttp

import (
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// TellHostPort 获得主机和端口，如果不存在则获取默认值.
func TellHostPort(host string, ssl bool) (server, port string, err error) {
	// 根据字符串获得IP和端口
	server, port, err = net.SplitHostPort(host)
	if err != nil {
		if addrerr, ok := err.(*net.AddrError); ok && strings.Contains(addrerr.Err, "missing port") {
			server = host
			if ssl {
				port = "443"
			} else {
				port = "80"
			}
			err = nil
		}
	}
	return server, port, err
}

// CombineURL 根据 path 构造 url
func CombineURL(r *http.Request, path string) *url.URL {
	return &url.URL{
		Scheme: r.URL.Scheme,
		Host:   r.Host,
		Path:   path,
	}
}

// IsDomainName 验证域名
func IsDomainName(domain string) bool {
	RegExp := regexp.MustCompile(`^(([a-zA-Z]{1})|([a-zA-Z]{1}[a-zA-Z]{1})|([a-zA-Z]{1}[0-9]{1})|([0-9]{1}[a-zA-Z]{1})|([a-zA-Z0-9][a-zA-Z0-9-_]{1,61}[a-zA-Z0-9]))\.([a-zA-Z]{2,6}|[a-zA-Z0-9-]{2,30}\.[a-zA-Z]{2,3})$`)

	return RegExp.MatchString(domain)
}
