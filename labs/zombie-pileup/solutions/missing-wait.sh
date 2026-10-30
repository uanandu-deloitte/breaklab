#!/usr/bin/env bash
set -euo pipefail
# Reference fix: reap each child as soon as it exits, instead of never
cat > /usr/local/bin/supervisor.py <<'PY'
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
        for c in list(children):
            try:
                done, _ = os.waitpid(c, os.WNOHANG)
                if done:
                    children.remove(c)
            except ChildProcessError:
                children.remove(c)
        time.sleep(0.3)


if __name__ == "__main__":
    main()
PY
restart-workers
