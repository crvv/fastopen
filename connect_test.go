package fastopen

import (
	"io"
	"net"
	"testing"
)

func echoTest(t *testing.T, conn net.Conn, message string) {
	_, err := conn.Write([]byte(message))
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()
	buffer := make([]byte, len(message))
	io.ReadFull(conn, buffer)
	if err != nil {
		t.Error(err)
	}
	if string(buffer) != message {
		t.Errorf("%v != %v", string(buffer), message)
	}
}
func TestDial(t *testing.T) {
	message := "tcp fastopen"
	done := make(chan int)
	listener, err := net.Listen("tcp", "[::]:12358")
	if err != nil {
		t.Error(err)
	}
	defer listener.Close()
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Error(err)
		}
		echoTest(t, conn, message)
		done <- 1
	}()
	conn, err := Dial("tcp-tfo", "localhost:12358")
	if err != nil {
		t.Error(err)
	}
	echoTest(t, conn, message)
	<-done
}
