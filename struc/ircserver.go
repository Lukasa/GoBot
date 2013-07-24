package struc

import (
	"net"
	"strconv"
	"strings"
)

/*
Represents a single IRC server. Currently there should be only one of these
in a running instance of GoBot, or it will get very confused.
*/
type IRCServer struct {
	IPAddr net.IP
	Port   uint64
	Name   string
}

// NewIRCServerFromHostnamePort creates a new IRCServer structure given a
// hostname-and-port string.
func NewIRCServerFromHostnamePort(hostnameAndPort string) (*IRCServer, error) {
	host, port := "", ""
	var err error

	if strings.Contains(hostnameAndPort, ":") {
		host, port, err = net.SplitHostPort(hostnameAndPort)
		if err != nil {
			return nil, err
		}
	} else {
		host = hostnameAndPort
		port = "6667"
	}

	numPort, err := strconv.ParseUint(port, 10, 64)
	if err != nil {
		return nil, err
	}

	addrs, err := net.LookupIP(host)
	if err != nil {
		return nil, err
	}

	// addrs can contain IPv6 addresses. For the moment, skip them. This is a bit of a dirty hack,
	// but it's easier than writing our own test function.
	var addr net.IP
	for _, address := range addrs {
		if address.To4() != nil {
			addr = address
			break
		}
	}

	server := new(IRCServer)
	server.IPAddr = addr
	server.Port = numPort
	server.Name = host

	return server, nil
}
