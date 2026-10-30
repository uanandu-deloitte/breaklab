# A container's public access level is what gates anonymous reads

Every Azure Blob Storage container has a `publicAccess` setting:
`off` (default — every request needs the account key, a SAS token, or
Azure AD auth), `blob` (anonymous read of individual blobs by URL, but
no listing), or `container` (anonymous read and listing). A blob that's
real and perfectly readable with the account key can still be
completely unreachable anonymously if the container's access level
doesn't allow it — that's not a bug in the blob, it's the container's
own setting.

## Why this matters

"It works when I use the connection string" and "it works anonymously"
are two different, independently-configured things. Teams regularly
ship a container meant to back a public status page or CDN origin, and
leave `publicAccess` at its secure-by-default `off` — the content is
real, the upload succeeded, and it's still unreachable by anyone who
doesn't already have a credential, which defeats the entire point of
publishing it publicly.

## Where to look

- `az storage container show --name public-assets --query properties.publicAccess`
  — the real, current access level.
- `curl -i http://127.0.0.1:10000/devstoreaccount1/public-assets/notice.txt`
  — a genuine anonymous request, no auth headers at all.

## The fix

```bash
az storage container set-permission --name public-assets --public-access blob
```

`blob` allows anonymous read of a blob by its direct URL without
exposing the ability to list everything in the container — the right
level for "this one file is meant to be public," not `container`.
