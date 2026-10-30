#!/usr/bin/env python3
"""Fixture endpoints for net-007. Three local ports are started here;
two more failure classes need no server code at all:

  9002  nothing listens here at all      -> connection refused
  9003  accepts the connection, then      -> timeout
        never reads or writes anything
  9004  a real HTTP server that returns  -> 404 for every path
        404 for every path
  9005  a real HTTP server that returns  -> 500 for every path
        500 for every path

(port 9001, used for the DNS-failure endpoint, is never bound by
anything here — the failure happens before any connection is even
attempted, because the hostname itself doesn't resolve.)
"""
import http.server
import socket
import socketserver
import threading


class NotFoundHandler(http.server.BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(404)
        self.end_headers()
        self.wfile.write(b"404: not found\n")

    def log_message(self, format, *args):
        pass


class ServerErrorHandler(http.server.BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(500)
        self.end_headers()
        self.wfile.write(b"500: internal server error\n")

    def log_message(self, format, *args):
        pass


class ReusableTCPServer(socketserver.TCPServer):
    allow_reuse_address = True


def serve_http(port, handler_cls):
    with ReusableTCPServer(("0.0.0.0", port), handler_cls) as httpd:
        httpd.serve_forever()


def serve_blackhole(port):
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    s.bind(("0.0.0.0", port))
    s.listen(5)
    while True:
        conn, _ = s.accept()
        # Deliberately never read from or write to conn — the client
        # just hangs here until its own timeout gives up. We also
        # never close it, so it stays open for the life of the probe.
        del conn


threading.Thread(target=serve_blackhole, args=(9003,), daemon=True).start()
threading.Thread(target=serve_http, args=(9004, NotFoundHandler), daemon=True).start()
threading.Thread(target=serve_http, args=(9005, ServerErrorHandler), daemon=True).start()

threading.Event().wait()
