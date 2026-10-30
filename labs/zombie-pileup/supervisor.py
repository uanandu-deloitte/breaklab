#!/usr/bin/env python3
import os
import sys
import time

PIDFILE = "/var/run/supervisor.pid"
REAPFLAG = "/tmp/breaklab-reap-flag"


def main():
    with open(PIDFILE, "w") as f:
        f.write(str(os.getpid()))
    children = []
    while True:
        if os.path.exists(REAPFLAG):
            for pid in children:
                try:
                    os.waitpid(pid, 0)
                except ChildProcessError:
                    pass
            sys.exit(0)
        pid = os.fork()
        if pid == 0:
            os._exit(0)
        children.append(pid)
        # BUG: completed children are never reaped here, so each one
        # stays a zombie until this process exits.
        time.sleep(0.3)


if __name__ == "__main__":
    main()
