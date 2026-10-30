#!/usr/bin/env python3
# Stand-in for the real service. It only ever writes its "graceful
# shutdown" marker from inside its own SIGTERM handler -- so the marker
# only appears if this process itself actually receives the signal.
import signal
import sys
import time

MARKER = "/var/run/app/graceful-shutdown.marker"


def handle_term(signum, frame):
    with open(MARKER, "w") as f:
        f.write("graceful\n")
    sys.exit(0)


signal.signal(signal.SIGTERM, handle_term)

while True:
    time.sleep(1)
