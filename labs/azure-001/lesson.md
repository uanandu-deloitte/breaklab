# Bicep type-checks parameters at compile time, before any deployment

`az bicep build` compiles a `.bicep` file into a plain ARM JSON
template — entirely locally, no Azure login or subscription needed.
Along the way it type-checks every parameter, resource property, and
expression against the real resource type schemas it ships with. A
parameter declared `int` with a quoted string literal as its default
isn't a style nitpick — Bicep rejects it outright, the same way any
statically-typed language would.

## Why this matters

Catching this at `bicep build` time — before anything is ever deployed
— is the entire point of using an IaC compiler instead of hand-writing
ARM JSON. A type error here would otherwise surface as a deployment
failure, or worse, a subtly wrong resource property, much later in a
pipeline.

## Where to look

- `main.bicep` — `param retentionDays int = '30'` declares the
  parameter as `int` but assigns it a string literal.
- `az bicep build --file main.bicep --stdout 2>&1` — the real compiler
  error names the exact line and the exact type mismatch.

## The fix

```bicep
param retentionDays int = 30
```

Drop the quotes — the value itself was always meant to be a number.
