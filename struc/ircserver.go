package struc

import (
	"net"
	"strconv"
	"strings"
)

/* Represents a single IRC server. Currently there should be only one of these
 * in a running instance of GoBot, or it will get very confused.
 */
type IRCServer struct {
	IPAddr net.IP
	Port   uint64
	Name   string
}

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

	// Return the first address only, because we're very lazy.
	addr := addrs[0]

	server := new(IRCServer)
	server.IPAddr = addr
	server.Port = numPort
	server.Name = host

	return server, nil
}
