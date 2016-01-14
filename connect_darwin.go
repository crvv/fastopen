package fastopen

/*
#include <sys/socket.h>
#include <netdb.h>
#include <strings.h>
#include <stdio.h>

int fastopenConnect(char* hostname, char* port, int family) {
	struct addrinfo *server;

	struct addrinfo hint;
	bzero((char*)&hint, sizeof(struct addrinfo));
	hint.ai_family = family;
	hint.ai_protocol = IPPROTO_TCP;
	hint.ai_socktype = SOCK_STREAM;

	int info = getaddrinfo(hostname, port, &hint, &server);
	if (info != 0) {
		return -5;
	}

	int fd = socket(server->ai_family, server->ai_socktype, server->ai_protocol);
	if (fd < 0) {
		freeaddrinfo(server);
		return -4;
	}

	sa_endpoints_t endpoints;
	bzero((char*)&endpoints, sizeof(endpoints));
	endpoints.sae_dstaddr = server->ai_addr;
	endpoints.sae_dstaddrlen = server->ai_addrlen;

	int rc = connectx(fd, &endpoints, SAE_ASSOCID_ANY,
		CONNECT_RESUME_ON_READ_WRITE | CONNECT_DATA_IDEMPOTENT,
		NULL, 0, NULL, NULL);
	if (rc < 0) {
		freeaddrinfo(server);
		return -3;
	}
	freeaddrinfo(server);
	return fd;
}
*/
import "C"

import (
	"fmt"
	"golang.org/x/sys/unix"
	"net"
	"time"
	"io"
)

type conn struct {
	fd int
}

func (c conn) Read(b []byte) (n int, err error) {
	n, err = unix.Read(c.fd, b)
	if n == 0 && err == nil {
		err = io.EOF
	}
	return
}
func (c conn) Write(b []byte) (n int, err error) {
	n, err = unix.Write(c.fd, b)
	return
}
func (c conn) Close() error {
	return unix.Close(c.fd)
}
func (c conn) LocalAddr() net.Addr {
	sa, err := unix.Getsockname(c.fd)
	if err != nil {
		return nil
	}
	return sockaddr2tcpaddr(sa)
}
func (c conn) RemoteAddr() net.Addr {
	sa, err := unix.Getpeername(c.fd)
	if err != nil {
		return nil
	}
	return sockaddr2tcpaddr(sa)
}
func (c conn) SetDeadline(t time.Time) (err error) {
	err = c.SetReadDeadline(t)
	if err != nil {
		return
	}
	err = c.SetWriteDeadline(t)
	return
}
func (c conn) SetReadDeadline(t time.Time) error {
	return unix.SetsockoptTimeval(c.fd, unix.SOL_SOCKET, unix.SO_RCVTIMEO, getTimeval(t))
}
func (c conn) SetWriteDeadline(t time.Time) error {
	return unix.SetsockoptTimeval(c.fd, unix.SOL_SOCKET, unix.SO_SNDTIMEO, getTimeval(t))
}

func getTimeval(t time.Time) *unix.Timeval {
	nsec := int64(t.Sub(time.Now()))
	timeval := unix.NsecToTimeval(nsec)
	return &timeval
}

func sockaddr2tcpaddr(addr unix.Sockaddr) (tcpaddr *net.TCPAddr) {
	switch a := addr.(type) {
	case *unix.SockaddrInet4:
		tcpaddr.Port = a.Port
		tcpaddr.IP = a.Addr[:]
	case *unix.SockaddrInet6:
		tcpaddr.Port = a.Port
		tcpaddr.IP = a.Addr[:]
		tcpaddr.Zone = fmt.Sprintf("%v", a.ZoneId)
	}
	return
}

func connect(network, address string) (c net.Conn, err error) {
	var family C.int
	switch network[len(network)-1] {
	case '4':
		family = C.AF_INET
	case '6':
		family = C.AF_INET6
	}
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return
	}
	fd := C.fastopenConnect(C.CString(host), C.CString(port), family)
	switch {
	case fd == -5:
		err = ErrorGetaddrinfo
	case fd == -4:
		err = ErrorSocket
	case fd == -3:
		err = ErrorConnectx
	case fd <= 0:
		err = UnknownError
	default:
		c = conn{fd: int(fd)}
	}
	return
}
