package fastopen

import "errors"

var NotSupport = errors.New("not supported")
var ErrorGetaddrinfo = errors.New("error in getaddrinfo")
var ErrorSocket = errors.New("error in socket")
var ErrorConnectx = errors.New("error in connectx")
var UnknownError = errors.New("unknown error")
