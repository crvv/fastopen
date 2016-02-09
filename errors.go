package fastopen

import (
	"errors"
	"golang.org/x/sys/unix"
)

var NotSupport = Err{errors.New("not supported")}
var ErrorGetaddrinfo = Err{errors.New("error in getaddrinfo")}
var ErrorSocket = Err{errors.New("error in socket")}
var ErrorConnectx = Err{errors.New("error in connectx")}
var UnknownError = Err{errors.New("unknown error")}

type Err struct {
	error
}

func (e Err) Timeout() bool {
	return e.error == unix.EAGAIN
}
func (e Err) Temporary() bool {
	return false
}
