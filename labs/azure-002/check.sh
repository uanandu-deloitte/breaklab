#!/usr/bin/env bash
set -uo pipefail

fail() { echo "FAIL: $1"; exit 1; }

for i in $(seq 1 30); do
  curl -sf -o /dev/null http://127.0.0.1:10000/devstoreaccount1 2>/dev/null
  [ $? -ne 7 ] && break
  sleep 1
done

code=$(curl -s -o /tmp/body.txt -w '%{http_code}' http://127.0.0.1:10000/devstoreaccount1/public-assets/notice.txt)
[ "$code" = "200" ] || fail "anonymous GET of public-assets/notice.txt returned HTTP $code, expected 200 -- the container still isn't readable without auth"

body=$(cat /tmp/body.txt)
expected="Checkout is back online. Thanks for your patience."
[ "$body" = "$expected" ] || fail "anonymous GET succeeded but returned unexpected content: $body"

access=$(az storage container show --name public-assets --query "properties.publicAccess" -o tsv)
[ "$access" = "blob" ] || fail "public-assets' public access level is \"$access\" -- blob-level anonymous read requires exactly \"blob\" (or \"container\"), not an authenticated workaround"

echo "PASS: public-assets/notice.txt is readable by a real anonymous HTTP GET, no auth, no SAS token"
