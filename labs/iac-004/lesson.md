# One template, many environments

Real infrastructure-as-code almost never hardcodes a value like an
environment name or a replica count directly into the config it
applies. Instead it keeps one **template** with placeholders, and a
separate small **variables file** per environment — `dev`, `staging`,
`prod` — that fills those placeholders in differently.

## Why this matters

Hardcoding `prod`-shaped values into your only copy of a config means
either keeping separate near-duplicate configs per environment (which
drift out of sync over time) or manually editing the one copy every
time you target a different environment (error-prone, and invisible
in a diff). Separating template from variables fixes both: the
template is shared and reviewed once, and each environment's
differences live in their own small, readable file.

## Where to look

This lab's apply tool is `infra-apply`. It takes an environment name as
its argument, loads the matching file under `/infra/vars/`, and
substitutes those variables into `/infra/resources.yaml.tmpl` before
writing the result to the real resource. Applying the same template
with a different environment argument produces a differently-shaped
resource — same placeholders, different values.
