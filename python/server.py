#!/usr/bin/env python3

# Echo server program
import socket

HOST = ''                 # Symbolic name meaning all available interfaces
PORT = 50007              # Arbitrary non-privileged port
s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.bind((HOST, PORT))
s.listen()
s.setsockopt(socket.SOL_TCP, 0x105, 1)

while True:
    conn, addr = s.accept()
    conn.sendall(b'tcp fastopen server')

    print('Connected by', addr)
    while True:
        data = conn.recv(1024)
        if not data:
            break
        print(data)
    conn.close()


