#!/usr/bin/env python3
# Stand-in for the real service: emits one heartbeat line per second to
# wherever /etc/app/logging.conf says logs should go.
import sys
import time


def log_target():
    try:
        with open("/etc/app/logging.conf") as f:
            for line in f:
                line = line.strip()
                if line.startswith("log_target="):
                    return line.split("=", 1)[1].strip()
    except FileNotFoundError:
        pass
    return "stdout"


target = log_target()
out = sys.stdout if target == "stdout" else open(target, "a", buffering=1)

n = 0
while True:
    n += 1
    print(f"heartbeat {n}", file=out, flush=True)
    time.sleep(1)
