#!/usr/bin/env bash
# Simulates sustained incoming job traffic.
while true; do
  touch "/var/spool/jobs/job-$RANDOM.job"
  sleep 0.03
done
