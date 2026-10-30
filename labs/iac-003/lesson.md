# Idempotency: safe to run twice

Every real infrastructure-as-code tool — Terraform, Ansible, Pulumi,
Chef — promises the same property: running `apply` again when nothing
has changed does **nothing**. That property is called **idempotency**,
and it's what makes these tools safe to re-run in CI on every commit,
on a schedule, or by hand after a worried "did that actually work?"

## Why this matters

Without idempotency, every re-run would blindly redo work: rewriting
files that already hold the right content, restarting services that
are already fine, recreating resources that already exist. Worse, a
tool that *isn't* idempotent makes "just run it again" a dangerous
habit instead of a safe one — which is exactly the habit IaC is
supposed to encourage.

A well-behaved apply tool instead does one simple thing on every run:
compare what's declared against what's really there, and act only on
the difference.

- **Nothing declared changed, reality already matches** → report "no
  changes" (or "up to date") and touch nothing.
- **Something declared changed, reality is now stale** → make the one
  change needed, and say so clearly.

## Where to look

This lab's apply tool is `infra-apply`. It reads a single declared
value from `/infra/resources.yaml` and makes sure the real resource
file matches it, printing a different message depending on whether it
actually had to change anything. Run it with nothing edited first, so
you see the no-op message before you ever see the "changed" one.
