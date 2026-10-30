#!/usr/bin/env python3
import http.server

METRICS = (
    "# HELP http_requests_count_total Total HTTP requests\n"
    "# TYPE http_requests_count_total counter\n"
    "http_requests_count_total 4321\n"
)


class Handler(http.server.BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path == "/metrics":
            body = METRICS.encode()
            self.send_response(200)
            self.send_header("Content-Type", "text/plain")
            self.send_header("Content-Length", str(len(body)))
            self.end_headers()
            self.wfile.write(body)
        else:
            self.send_response(404)
            self.end_headers()

    def log_message(self, fmt, *args):
        pass


if __name__ == "__main__":
    http.server.HTTPServer(("0.0.0.0", 9091), Handler).serve_forever()
