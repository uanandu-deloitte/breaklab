#!/usr/bin/env python3
"""The "database" for net-900: a tiny local TCP service standing in for
a real database without needing one. It only ever does one thing --
reply DB_OK to anything that connects -- so the api process's own
request to it is a genuine second network hop, not a stub it fakes
internally."""
import socketserver


class Handler(socketserver.BaseRequestHandler):
    def handle(self):
        self.request.recv(1024)
        self.request.sendall(b"DB_OK\n")


class ReusableTCPServer(socketserver.TCPServer):
    allow_reuse_address = True


with ReusableTCPServer(("127.0.0.1", 6000), Handler) as srv:
    srv.serve_forever()
