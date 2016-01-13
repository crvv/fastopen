package fastopen

import (
	"net"
	"strings"
)

func Dial(network, address string) (c net.Conn, err error) {
	switch network {
	case "tcp", "tcp4", "tcp6":
		c, err = net.Dial(network, address)
		return
	case "tcp-tfo", "tcp4-tfo", "tcp6-tfo":
		c, err = connect(strings.TrimRight(network, "-tfo"), address)
		return
	default:
		err = NotSupport
		return
	}
}
