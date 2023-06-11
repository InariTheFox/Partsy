package utils

import (
	"net"
	"net/http"

	"github.com/inarithefox/partsy/server/public/model"
)

func GetIPAddress(r *http.Request) string {
	host, _, _ := net.SplitHostPort(r.RemoteAddr)

	return host
}

func GetProtocol(r *http.Request) string {
	if r.Header.Get(model.HeaderForwardedProto) == "https" || r.TLS != nil {
		return "https"
	}

	return "http"
}
