#!/usr/bin/env python3
"""Backend API for net-900. Reads its own config file at startup (exits
if it can't -- same idea as lnx-015's config-permission fault), then on
every request makes a real local TCP call to the "database" at
db.internal:6000 and only answers 200 if that call actually succeeds."""
import http.server
import socket
import socketserver
import sys

CONFIG = "/opt/api/config.yaml"
DB_HOST = "db.internal"
DB_PORT = 6000

try:
    with open(CONFIG) as f:
        f.read()
except OSError as e:
    sys.stderr.write("cannot read config: {}\n".format(e))
    sys.exit(1)


class Handler(http.server.BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path != "/":
            self.send_response(404)
            self.end_headers()
            return

        try:
            with socket.create_connection((DB_HOST, DB_PORT), timeout=2) as s:
                s.settimeout(2)
                s.sendall(b"ping")
                reply = s.recv(1024)
        except OSError as e:
            self.send_response(500)
            self.end_headers()
            self.wfile.write(("FAIL: database unreachable: " + str(e) + "\n").encode())
            return

        if reply.strip() != b"DB_OK":
            self.send_response(500)
            self.end_headers()
            self.wfile.write(b"FAIL: unexpected database response\n")
            return

        self.send_response(200)
        self.end_headers()
        self.wfile.write(b"BreakLab OK: net900-api-db-chain-healthy\n")

    def log_message(self, format, *args):
        pass


class ReusableTCPServer(socketserver.TCPServer):
    allow_reuse_address = True


with ReusableTCPServer(("127.0.0.1", 5000), Handler) as httpd:
    httpd.serve_forever()
