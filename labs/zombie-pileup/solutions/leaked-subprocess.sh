#!/usr/bin/env bash
set -euo pipefail
# Reference fix: wait on each Popen job once it's done, instead of never
cat > /usr/local/bin/worker_pool.py <<'PY'
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
        for proc in list(procs):
            if proc.poll() is not None:
                procs.remove(proc)
        time.sleep(0.3)


if __name__ == "__main__":
    main()
PY
restart-workers
