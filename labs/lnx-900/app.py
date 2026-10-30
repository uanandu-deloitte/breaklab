#!/usr/bin/env python3
"""Backend application for lnx-900: reads its config file at startup
(exits if it can't), and on every request checks how much space its
local data directory is using — standing in for a real "disk full"
condition without needing to actually exhaust the container's real
disk."""
import http.server
import os
import socketserver
import sys
import time

CONFIG = "/opt/app/config.yaml"
DATA_DIR = "/var/log/app-data"
FULL_THRESHOLD_MB = 10

try:
    with open(CONFIG) as f:
        f.read()
except OSError as e:
    sys.stderr.write("cannot read config: {}\n".format(e))
    sys.exit(1)


def dir_size_mb(path):
    total = 0
    for root, _dirs, files in os.walk(path):
        for name in files:
            try:
                total += os.path.getsize(os.path.join(root, name))
            except OSError:
                pass
    return total / (1024 * 1024)


class Handler(http.server.BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path != "/":
            self.send_response(404)
            self.end_headers()
            return

        size_mb = dir_size_mb(DATA_DIR)
        if size_mb >= FULL_THRESHOLD_MB:
            self.send_response(500)
            self.end_headers()
            self.wfile.write(b"FAIL: no space left on backend data volume\n")
            return

        try:
            with open(os.path.join(DATA_DIR, "access.log"), "a") as f:
                f.write("{} request\n".format(time.time()))
        except OSError:
            pass

        self.send_response(200)
        self.end_headers()
        self.wfile.write(b"BreakLab OK\n")

    def log_message(self, format, *args):
        pass


class ReusableTCPServer(socketserver.TCPServer):
    allow_reuse_address = True


with ReusableTCPServer(("127.0.0.1", 5000), Handler) as httpd:
    httpd.serve_forever()
