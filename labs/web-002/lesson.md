# Reverse proxies, in one directive

A reverse proxy sits in front of one or more backend servers and
forwards client requests to them, returning the backend's response as
if the proxy had produced it itself. nginx does this with a single
directive: `proxy_pass`.

## The minimum config

```
location / {
    proxy_pass http://127.0.0.1:8090;
}
```

Any request matching `location /` is now forwarded to whatever is
listening on `127.0.0.1:8090`, and nginx relays the response back to
the client. Without this directive, nginx just serves its own
configured document root — which is why an un-proxied nginx shows its
own default page no matter what's running elsewhere.

## Headers matter too

By default, a proxied request's `Host` header still carries the
address nginx itself is listening on, not the one the client actually
asked for. Most backends that care about virtual hosting, URL
generation, or logging want the original value:

```
proxy_set_header Host $host;
```

`$host` is an nginx variable holding the request's original `Host`
header (or the server name, if none was sent) — forwarding it keeps the
backend's view of the request consistent with what the client sent.

## Validate before you reload

`nginx -t` checks a config's syntax without touching the running
process — always run it before `nginx -s reload`. A config with a typo
fails `-t` immediately, before it can take down a server that was
previously working fine.
