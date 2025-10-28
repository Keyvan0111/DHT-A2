package utils

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
)

var domainSuffix = func() string {
	val := strings.Trim(strings.TrimSpace(os.Getenv("DOMAIN_SUFFIX")), ".")
	if val == "" {
		return "ifi.uit.no"
	}
	return val
}()

func useDomain(host string) bool {
	if host == "" {
		return false
	}
	if strings.EqualFold(host, "localhost") {
		return false
	}
	if net.ParseIP(host) != nil {
		return false
	}
	return !strings.Contains(host, ".")
}

func FullHost(host string) string {
	h := strings.TrimSpace(host)
	if h == "" {
		return ""
	}
	if useDomain(h) {
		return fmt.Sprintf("%s.%s", h, domainSuffix)
	}
	return h
}

func BuildHTTPAddr(host, port string) string {
	return fmt.Sprintf("http://%s:%s", FullHost(host), port)
}

func NormalizeHostPort(input string) (string, string, error) {
	raw := strings.TrimSpace(input)
	if raw == "" {
		return "", "", errors.New("empty address")
	}

	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		u, err := url.Parse(raw)
		if err != nil {
			return "", "", err
		}
		host := u.Hostname()
		port := u.Port()
		if port == "" {
			return "", "", errors.New("missing port")
		}
		return host, port, nil
	}

	if !strings.Contains(raw, ":") {
		return "", "", errors.New("address must be HOST:PORT")
	}

	host, port, found := strings.Cut(raw, ":")
	if !found || host == "" || port == "" {
		return "", "", errors.New("malformed HOST:PORT pair")
	}

	return host, port, nil
}
