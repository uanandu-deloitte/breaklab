# Disk Full Incident

Severity: High. The application stopped writing logs and new requests
fail intermittently. Restore safe operation without deleting protected
customer data.

## Tasks

- find what is consuming space under `/var/log`
- remove ONLY the unnecessary generated/junk data
- confirm the protected fixtures are untouched:
  - `/srv/data/customer.db`
  - `/opt/app/bin/app`
  - `/var/log/audit/retained.log`
- confirm the application can write new logs again (e.g. a new file under
  `/var/log/app/`)

## Constraints

- do not touch `/srv/data/customer.db`
- do not touch `/opt/app/bin/app`
- do not touch `/var/log/audit/retained.log`

Get a shell: `docker exec -it breaklab-lnx-022 bash`
