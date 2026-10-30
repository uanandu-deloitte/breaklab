# Idempotency is a property of the module, not the intent

Every real configuration-management tool — Ansible, Terraform, Chef,
Puppet — promises the same thing: running it again when nothing needs
to change should do **nothing**. Ansible's purpose-built modules
(`lineinfile`, `copy`, `template`, `package`, ...) check real state
before acting and report `changed` only when they actually changed
something. `shell:` and `command:` tasks don't — they just run, every
time, with no idea whether they needed to.

## Why this matters

`shell: echo "x" >> file` is not "set x in file" — it's "append x to
file," unconditionally, forever. The first run looks like it worked.
The second run proves it didn't actually understand the goal: the file
now has two lines, and a hundred runs later it has a hundred.

## The fix

Replace the `shell:`/`command:` task with the real module that models
the actual intent — here, "this exact line should exist, once":

```yaml
- name: ensure max_connections is configured
  ansible.builtin.lineinfile:
    path: /etc/app/app.conf
    regexp: '^max_connections='
    line: 'max_connections=200'
```

`lineinfile` checks the file first: if the line's already correct, it
reports `ok` with no change. Run it a second time and `ansible-playbook`
should report `changed=0`.
