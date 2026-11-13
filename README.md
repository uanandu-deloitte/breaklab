# BreakLab

A break/fix DevOps training ground: break a throwaway Docker environment at
random, hand the learner a symptom-only incident ticket, make them fix it
with no instructions, and keep a postmortem they own.

- **Outcome-based, not solution-based.** `breaklab check` verifies the
  symptom is gone, not that you ran a specific command — any valid fix
  passes.
- **The write-up is part of the lab.** A lab isn't done until its
  postmortem (Summary, Detection, Root cause, Fix, Prevention) has real
  content.
- **One static binary.** Go + the `docker` CLI are the only runtime
  dependencies; everything else (labs, roadmap, progress) is plain files.

## Quick start

```bash
git clone <this-repo> breaklab && cd breaklab

cd frontend && npm ci && npm run build && cd ..   # builds frontend/dist/
go build -o breaklab ./cmd/breaklab

./breaklab list                          # see the catalog
./breaklab start web-down                # break it
docker exec -it breaklab-web-down bash   # go fix it
./breaklab check                         # did it work?
./breaklab postmortem                    # write up what happened
./breaklab serve                         # dashboard at http://127.0.0.1:8787
```

## Documentation

Pick the doc for the hat you're wearing — each assumes you've read none of
the others.

| If you are... | Read | For |
|---|---|---|
| Building or extending BreakLab | [`docs/developer-guide.md`](docs/developer-guide.md) | Architecture, CLI internals, lab format, testing, contributing |
| A learner working through labs | [`docs/learner-guide.md`](docs/learner-guide.md) | The loop, the phases, how to use hints and postmortems honestly |
| Provisioning learner sandboxes in the cloud | [`docs/devops-guide-azure.md`](docs/devops-guide-azure.md) | Standing up isolated, disposable instances in Azure (GCP/AWS to follow) |
| Deploying one instance on your own LAN | [`docs/deploy-raspberry-pi.md`](docs/deploy-raspberry-pi.md) | A single headless instance (e.g. Raspberry Pi 5) for personal use |
| Reviewing the original build spec | [`docs/breaklab-spec.md`](docs/breaklab-spec.md) | The full spec-driven build plan |

## Project status

150+ labs span the curriculum phases from Linux fundamentals to a
capstone — see `docs/developer-guide.md` for the current count and
layout.

## References

- [devopslabs.dev](https://devopslabs.dev) — the hands-on platform this
  project takes its format from, reimplemented standalone.

## License

[MIT](LICENSE)
