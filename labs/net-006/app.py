#!/usr/bin/env python3
"""Fixture API for net-006. A request is only accepted when ALL of the
following are true:

  - method is POST
  - header Content-Type: application/json is present
  - header Authorization: Bearer test-token-fixture is present
    (a synthetic fixture token, never a real credential)
  - the JSON body contains {"action": "sync"}

Every request's outcome — success or the specific reason it failed —
is persisted to RESULT_FILE so check.sh can confirm the *learner's own*
curl call actually succeeded, instead of check.sh re-sending the
request itself."""
import http.server
import json
import socketserver

RESULT_FILE = "/tmp/last-request-result.txt"
REQUIRED_AUTH = "Bearer test-token-fixture"


def record(outcome):
    with open(RESULT_FILE, "w") as f:
        f.write(outcome + "\n")


class Handler(http.server.BaseHTTPRequestHandler):
    def _reply(self, status, outcome, body):
        record(outcome)
        self.send_response(status)
        self.send_header("Content-Type", "application/json")
        self.end_headers()
        self.wfile.write(body)

    def _handle(self, method):
        if method != "POST":
            self._reply(400, f"fail: wrong method {method}",
                        b'{"error": "use POST"}\n')
            return

        content_type = self.headers.get("Content-Type", "")
        if content_type.split(";")[0].strip() != "application/json":
            self._reply(400, "fail: missing or wrong Content-Type",
                        b'{"error": "Content-Type must be application/json"}\n')
            return

        if self.headers.get("Authorization", "") != REQUIRED_AUTH:
            self._reply(401, "fail: missing or wrong Authorization",
                        b'{"error": "Authorization required"}\n')
            return

        length = int(self.headers.get("Content-Length", 0) or 0)
        raw = self.rfile.read(length) if length else b""
        try:
            payload = json.loads(raw) if raw else {}
        except json.JSONDecodeError:
            self._reply(400, "fail: invalid JSON body",
                        b'{"error": "invalid JSON body"}\n')
            return

        if payload.get("action") != "sync":
            self._reply(400, "fail: missing required field action=sync",
                        b'{"error": "body must include action: sync"}\n')
            return

        self._reply(200, "success", b'{"result": "accepted"}\n')

    def do_GET(self):
        self._handle("GET")

    def do_POST(self):
        self._handle("POST")

    def do_PUT(self):
        self._handle("PUT")

    def do_DELETE(self):
        self._handle("DELETE")

    def log_message(self, format, *args):
        pass


class ReusableTCPServer(socketserver.TCPServer):
    allow_reuse_address = True


with ReusableTCPServer(("0.0.0.0", 8080), Handler) as httpd:
    httpd.serve_forever()
