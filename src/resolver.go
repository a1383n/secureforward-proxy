package src

import (
	"context"
	"fmt"
	"net"
	"time"
)

func resolveIPAddress(hostname string) (string, error) {
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Second,
			}
			return d.DialContext(ctx, "udp", "8.8.8.8:53")
		},
	}

	addrs, err := resolver.LookupIP(context.Background(), "ip4", hostname)
	if err != nil {
		return "", err
	}

	if len(addrs) == 0 {
		return "", fmt.Errorf("no IP addresses found for host %s", hostname)
	}

	return addrs[0].String(), nil
}
