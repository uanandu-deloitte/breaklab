#!/usr/bin/env python3
"""Tiny health-check app for lnx-015: re-reads its config file on every
request (not just at startup) so a permission fix takes effect
immediately, with no process restart required."""
import http.server
import socketserver

CONFIG = "/srv/app/config.yaml"


class Handler(http.server.BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path != "/healthz":
            self.send_response(404)
            self.end_headers()
            return
        try:
            with open(CONFIG) as f:
                f.read()
        except OSError as e:
            self.send_response(500)
            self.end_headers()
            self.wfile.write(("FAIL: cannot read config: " + str(e) + "\n").encode())
            return
        self.send_response(200)
        self.end_headers()
        self.wfile.write(b"OK: config loaded\n")

    def log_message(self, format, *args):
        pass


class ReusableTCPServer(socketserver.TCPServer):
    allow_reuse_address = True


with ReusableTCPServer(("0.0.0.0", 80), Handler) as httpd:
    httpd.serve_forever()
