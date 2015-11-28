#!/usr/bin/env python3

# Echo client program
import socket
import fastopen

MSG_FASTOPEN = 0x20000000

HOST = 'localhost'    # The remote host
PORT = 50007              # The same port as used by the server
s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

fastopen.connect(s.fileno(), HOST, PORT)

s.send(b'Hello, world')

data = s.recv(1024)
s.close()
print('Received', repr(data))
