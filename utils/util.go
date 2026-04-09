package utils

import (
	"net"
	"net/http"
	"net/url"
	"strings"
)

func IsValidURL(raw string) bool {
	u, err := url.ParseRequestURI(raw)
	if err != nil {
		return false
	}

	// require http or https
	// if u.Scheme != "http" && u.Scheme != "https" {
	// 	return false
	// }

	// must have a host
	if u.Host == "" {
		return false
	}

	return true
}

func GetClientIP(r *http.Request) string {
	// Try X-Forwarded-For
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// Format: client, proxy1, proxy2
		parts := strings.Split(xff, ",")
		ip := strings.TrimSpace(parts[0])
		if ip != "" {
			return ip
		}
	}

	// Fallback to RemoteAddr
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}

	return r.RemoteAddr
}
