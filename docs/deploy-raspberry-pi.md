# Deploying BreakLab to a Raspberry Pi 5

BreakLab compiles to a single statically-linked Go binary, so running it
on a Pi 5 is mostly: install Docker, copy one file plus your data
directories, and start it. This guide assumes the Pi is headless (no
monitor attached) and that you'll browse the local web UI from your Mac
or another device on the same home network.

## 1. Prerequisites

- A Raspberry Pi 5 (4 GB or 8 GB), with Raspberry Pi OS (64-bit, Bookworm)
  or Ubuntu Server 24.04 LTS (arm64) installed and SSH enabled.
- The Pi and your dev machine on the same LAN.
- A free TCP port for the web UI (default `8787`).

Confirm the OS is 64-bit before continuing — BreakLab's `linux/arm64`
build will not run on a 32-bit userland:

```bash
ssh pi@<pi-host>.local uname -m
# expect: aarch64
```

## 2. Install Docker on the Pi

BreakLab shells out to the `docker` CLI, so Docker Engine must be
installed and running on the Pi itself (labs run as containers on the
Pi, not on your Mac):

```bash
ssh pi@<pi-host>.local
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker "$USER"
# log out and back in (or `newgrp docker`) for the group change to apply
docker run --rm hello-world   # sanity check
```

## 3. Build the binary and frontend for the Pi

Cross-compile from your Mac rather than building on the Pi — it's
faster and keeps the Pi free of a Go toolchain. Build the frontend here
too: it's architecture-independent static output (Node is build-time
only — the Pi never runs a Node process), so there's no cross-compile
step for it, just a regular `npm run build`:

```bash
cd breaklab
GOOS=linux GOARCH=arm64 go build -o dist/breaklab-arm64 ./cmd/breaklab
(cd frontend && npm ci && npm run build)   # produces frontend/dist/
```

Copy the binary and the data BreakLab needs onto the Pi. The binary and
`frontend/dist/` are the only built artifacts; `labs/`, `templates/`,
and `roadmap.yaml` are plain files it reads at runtime. `--frontend`
defaults to `frontend/dist` relative to the working directory
`breaklab-arm64` is started from, so keep that same relative layout on
the Pi:

```bash
rsync -avz dist/breaklab-arm64 labs/ templates/ roadmap.yaml \
  pi@<pi-host>.local:~/breaklab/
ssh pi@<pi-host>.local mkdir -p ~/breaklab/frontend/dist
rsync -avz frontend/dist/ pi@<pi-host>.local:~/breaklab/frontend/dist/
```

If you'd rather build on-device (e.g. no cross-compile toolchain handy),
install Go on the Pi (`sudo apt install golang-go`, or a newer version
from go.dev if the distro package is stale) and run `go build` there
instead — same output, just slower.

## 4. First run

```bash
ssh pi@<pi-host>.local
cd ~/breaklab
chmod +x breaklab-arm64
./breaklab-arm64 list          # confirms the catalog loads
./breaklab-arm64 serve --host 0.0.0.0 --port 8787
```

`--host 0.0.0.0` is required here — the default bind is `127.0.0.1`,
which on a headless Pi would only be reachable from a shell on the Pi
itself. The CLI prints a one-time warning when you pass a non-loopback
host: the UI has no authentication, so anything that can reach the Pi
on port 8787 can see it. Only do this on a trusted home network, never
on a network you don't control, and never port-forward 8787 to the
internet.

From your Mac, open `http://<pi-host>.local:8787` (or the Pi's LAN IP)
in a browser.

This Pi deployment is the unauthenticated, local/LAN trust model —
fine for a home network, never the internet. `breaklab serve --db
<path>` is a different, separate deployment shape: it adds real
accounts (passwords, sessions, API tokens) for `breaklab sync`, and
anything handling those belongs behind TLS with real exposure
planning, not just a one-time warning. Don't add `--db` to this
LAN-only setup without doing that first.

## 5. Run it as a systemd service

So it survives reboots and restarts on crash, install it as a service
instead of running it in a terminal:

```bash
sudo tee /etc/systemd/system/breaklab.service >/dev/null <<'EOF'
[Unit]
Description=BreakLab local web UI
After=docker.service
Requires=docker.service

[Service]
Type=simple
User=pi
WorkingDirectory=/home/pi/breaklab
ExecStart=/home/pi/breaklab/breaklab-arm64 serve --host 0.0.0.0 --port 8787
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable --now breaklab.service
systemctl status breaklab.service
```

Logs: `journalctl -u breaklab.service -f`.

## 6. Where your data lives

`progress/progress.json` and `postmortems/*.md` are written under
`~/breaklab/` on the Pi, same as on your Mac. If you want them backed up
and versioned, make `~/breaklab/` a git clone of the same repo and push
from the Pi, or periodically `rsync` the `progress/` and `postmortems/`
directories back to your Mac and commit there — either works, since
they're plain files either way. `.breaklab/session.json` stays
gitignored wherever it runs.

## 7. Updating after a new OpenSpec change

Each time you archive a change that touches the binary, `frontend/`, or
`templates/`, rebuild and redeploy:

```bash
GOOS=linux GOARCH=arm64 go build -o dist/breaklab-arm64 ./cmd/breaklab
(cd frontend && npm ci && npm run build)
rsync -avz dist/breaklab-arm64 templates/ pi@<pi-host>.local:~/breaklab/
rsync -avz frontend/dist/ pi@<pi-host>.local:~/breaklab/frontend/dist/
ssh pi@<pi-host>.local sudo systemctl restart breaklab.service
```

If a change only added or edited files under `labs/`, sync just that
directory and restart — the binary itself hasn't changed.

## 8. Troubleshooting

- **`permission denied` on `docker run`**: the `usermod -aG docker` group
  change from step 2 needs a fresh login session — log out/in or reboot.
- **Port already in use**: `breaklab serve` fails fast with the port
  number named in the error rather than silently picking another one;
  stop whatever else is bound to it or choose a different `--port`.
- **Can't reach `http://<pi-host>.local:8787` from another device**:
  confirm `breaklab serve` was started with `--host 0.0.0.0` (not the
  default `127.0.0.1`), and that no firewall (e.g. `ufw`) on the Pi is
  blocking the port.
- **`exec format error` running the binary**: you copied an
  `amd64`/`darwin` build instead of the `GOOS=linux GOARCH=arm64` one —
  rebuild with the correct `GOOS`/`GOARCH` pair from step 3.
