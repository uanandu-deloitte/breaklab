# Building a real request with curl

`curl` defaults to a bare `GET` with no extra headers and no body — fine
for fetching a page, useless for exercising an API that expects something
more specific. Reproducing a real API call means controlling four things
independently.

## Method

`-X <METHOD>` overrides curl's default `GET`. `-X POST` sends a POST
request. (Using `-d` alone also implies POST — but being explicit with
`-X` makes the intent unambiguous and is worth the habit.)

## Headers

`-H '<Name>: <value>'` adds one request header. Pass it multiple times
for multiple headers:

```bash
curl -X POST \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer <token>' \
  ...
```

A server that checks for a specific header checks for it exactly — a
missing header and a misspelled one fail the same way.

## Body

`-d '<data>'` sends request data as the body. Whether a server parses
that body as JSON depends entirely on the `Content-Type` header you set
alongside it — `-d` by itself doesn't declare a content type; pairing it
with `Content-Type: application/json` is what tells the server to expect
and parse JSON rather than a URL-encoded form.

## Reading what came back

A non-2xx status and a response body usually explain exactly what was
missing or wrong — `curl -i` (headers) or `curl -v` (full
request/response trace) surface it right away rather than leaving you to
guess.
