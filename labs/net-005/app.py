#!/usr/bin/env python3
"""Tiny backend fixture for net-005: a frontend would reach this over
the hostname api.internal, never a raw IP. It always answers the same
JSON marker so the check can prove the *name* resolved and the
request actually reached this process."""
import http.server
import socketserver

MARKER = b'{"service": "api-internal", "status": "ok"}\n'


class Handler(http.server.BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.send_header("Content-Type", "application/json")
        self.end_headers()
        self.wfile.write(MARKER)

    def log_message(self, format, *args):
        pass


class ReusableTCPServer(socketserver.TCPServer):
    allow_reuse_address = True


with ReusableTCPServer(("0.0.0.0", 8080), Handler) as httpd:
    httpd.serve_forever()
