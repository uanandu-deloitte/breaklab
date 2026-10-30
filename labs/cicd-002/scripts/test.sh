#!/usr/bin/env bash
# Fixture data: this test suite genuinely fails, on purpose, every time it
# runs. This lab is about how run.sh handles that failure, not about
# fixing test.sh itself — it must keep failing forever.
echo "Running test suite..."
echo "FAILED: 1 of 12 tests failed (test_payment_timeout)"
exit 1
