import subprocess, time, os, glob

CAP = 10
WORKER_SECONDS = 3
QUEUE = "/var/spool/jobs"
ACKDIR = "/var/spool/jobs/done"
STATUS = "/var/log/handler-status.log"


def log(msg):
    with open(STATUS, "a") as f:
        f.write(msg + "\n")


procs = []
while True:
    procs = [p for p in procs if p.poll() is None]
    slots = CAP - len(procs)
    if slots > 0:
        pending = sorted(glob.glob(QUEUE + "/*.job"))[:slots]
        for path in pending:
            jobid = os.path.basename(path)
            try:
                os.remove(path)
            except FileNotFoundError:
                continue
            try:
                cmd = "touch " + ACKDIR + "/" + jobid + " && sleep " + str(WORKER_SECONDS) + "; true"
                p = subprocess.Popen(["sh", "-c", cmd])
                procs.append(p)
            except OSError as e:
                log("ERROR: cannot start worker for " + jobid + ": " + str(e))
    time.sleep(0.2)
