# Solution: HTTP Failure Classification

```bash
curl -v --max-time 3 http://nosuchhost.invalid:9001/
# curl: (6) Could not resolve host: nosuchhost.invalid  -> dns-failure

curl -v --max-time 3 http://127.0.0.1:9002/
# curl: (7) Failed to connect ... Connection refused      -> refused

curl -v --max-time 3 http://127.0.0.1:9003/
# curl: (28) Operation timed out after 3000 milliseconds  -> timeout

curl -v --max-time 3 http://127.0.0.1:9004/anything
# < HTTP/1.0 404 Not Found                                -> 404

curl -v --max-time 3 http://127.0.0.1:9005/
# < HTTP/1.0 500 Internal Server Error                    -> 500

cat <<EOF > /workspace/answers.env
ENDPOINT_1=dns-failure
ENDPOINT_2=refused
ENDPOINT_3=timeout
ENDPOINT_4=404
ENDPOINT_5=500
EOF
```

`curl -v` labels the exact failure stage on its own: error 6 is a
resolver failure (DNS), error 7 is a refused TCP connection, error 28 is a
client-side timeout waiting for a response that never came. Once a
connection actually succeeds and a real HTTP response comes back, the
status code alone tells client error (404) apart from server error
(500) — both reached a working server, but they mean opposite things.
