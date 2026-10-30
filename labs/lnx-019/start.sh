#!/usr/bin/env bash
# start.sh - bring up the app
#
# Bug 1: APP_CONFIG_PATH is required but nothing checks it's actually
# set before it's used.
# Bug 2: the real status of the config check is thrown away by piping it
# into grep -q and then never looking at $? — the script below always
# reaches the final echo and always exits 0, no matter what happened.
cat "$APP_CONFIG_PATH" | grep -q "service:"
echo "Service started successfully"
exit 0
