#!/usr/bin/env python3
import os
import subprocess
import sys
import time

PIDFILE = "/var/run/worker_pool.pid"
REAPFLAG = "/tmp/breaklab-reap-flag"


def main():
    with open(PIDFILE, "w") as f:
        f.write(str(os.getpid()))
    procs = []
    while True:
        if os.path.exists(REAPFLAG):
            for p in procs:
                try:
                    p.wait(timeout=2)
                except Exception:
                    pass
            sys.exit(0)
        p = subprocess.Popen(["/bin/true"])
        procs.append(p)
        # BUG: finished jobs are never reaped with .wait()/.communicate(),
        # so each completed Popen leaves a zombie behind.
        time.sleep(0.3)


if __name__ == "__main__":
    main()
