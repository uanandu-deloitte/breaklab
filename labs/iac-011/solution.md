# Solution: The Playbook That's Never Actually Done

```yaml
---
- hosts: localhost
  connection: local
  tasks:
    - name: ensure max_connections is configured
      ansible.builtin.lineinfile:
        path: /etc/app/app.conf
        regexp: '^max_connections='
        line: 'max_connections=200'
```

```bash
cd /infra
ansible-playbook -i inventory playbook.yml
# PLAY RECAP: localhost : ok=1 changed=1 ...

ansible-playbook -i inventory playbook.yml
# PLAY RECAP: localhost : ok=1 changed=0 ...  <- idempotent
```
