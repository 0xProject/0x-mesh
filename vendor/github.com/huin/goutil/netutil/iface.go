package netutil

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
)

// ExtractAddrParts extracts parts out of an address that is one of
// *net.IPAddr, *net.TCPAddr, or *net.UDPAddr.
func ExtractAddrParts(addr net.Addr) (ip net.IP, port int, ip6Zone string, err error) {
	switch addr := addr.(type) {
	case *net.IPAddr:
		ip = addr.IP
		ip6Zone = addr.Zone
	case *net.TCPAddr:
		ip = addr.IP
		port = addr.Port
		ip6Zone = addr.Zone
	case *net.UDPAddr:
		ip = addr.IP
		port = addr.Port
		ip6Zone = addr.Zone
	default:
		err = errors.New("unknown address type")
		return
	}
	return
}

// FormAddr creates a net.Addr for the given network type.
func FormAddr(network string, ip net.IP, port int, ip6Zone string) (addr net.Addr, err error) {
	switch network {
	case "ip4", "tcp4", "udp4":
		if len(ip) != net.IPv4len {
			err = fmt.Errorf("invalid IP address %s for network %s", ip, network)
			return
		}
	case "ip6", "tcp6", "udp6":
		if len(ip) != net.IPv6len {
			err = fmt.Errorf("invalid IP address %s for network %s", ip, network)
			return
		}
	}

	switch network {
	case "ip", "ip4", "ip6":
		addr = &net.IPAddr{IP: ip, Zone: ip6Zone}
	case "tcp", "tcp4", "tcp6":
		addr = &net.TCPAddr{IP: ip, Port: port, Zone: ip6Zone}
	case "udp", "udp4", "udp6":
		addr = &net.UDPAddr{IP: ip, Port: port, Zone: ip6Zone}
	default:
		err = fmt.Errorf("unknown network type: %s", network)
		return
	}

	return
}

func IsZeroIP(ip net.IP) bool {
	switch len(ip) {
	case 0:
		return true
	case net.IPv4len:
		return bytes.Equal([]byte(ip), net.IPv4zero)
	case net.IPv6len:
		return bytes.Equal([]byte(ip), net.IPv6zero)
	}
	return false
}

func ExpandIP(ip net.IP) ([]net.IP, error) {
	if !IsZeroIP(ip) {
		return []net.IP{ip}, nil
	}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	ips := make([]net.IP, 0, len(addrs))
	for _, addr := range addrs {
		addr, ok := addr.(*net.IPNet)
		if !ok {
			// Ignore address - does not have an IP address.
			continue
		}
		ips = append(ips, addr.IP)
	}
	return ips, nil
}

func ExpandHostedAddr(addr net.Addr) (addresses []net.Addr, err error) {
	ip, port, ip6Zone, err := ExtractAddrParts(addr)
	if err != nil {
		return
	}

	network := addr.Network()

	ips, err := ExpandIP(ip)
	if err != nil {
		return
	}

	addresses = make([]net.Addr, 0, len(ips))
	for i := range ips {
		addr, addrErr := FormAddr(network, ips[i], port, ip6Zone)
		if addrErr != nil {
			continue
		}
		addresses = append(addresses, addr)
	}

	return
}

func ExpandHostedAddrToURLs(addr net.Addr, tmplURL url.URL) (urls []url.URL, err error) {
	var addresses []net.Addr
	if addresses, err = ExpandHostedAddr(addr); err != nil {
		return
	}

	urls = make([]url.URL, 0, len(addresses))
	for i := range addresses {
		ip, port, ip6Zone, extractErr := ExtractAddrParts(addresses[i])
		if extractErr != nil {
			err = extractErr
			return
		}

		var ipStr string
		if len(ip) == net.IPv6len && ip6Zone != "" {
			ipStr = fmt.Sprintf("%s%%%s", ip, ip6Zone)
		} else {
			ipStr = ip.String()
		}

		u := tmplURL
		u.Host = net.JoinHostPort(ipStr, strconv.Itoa(port))
		urls = append(urls, u)
	}

	return
}
