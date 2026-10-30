#!/usr/bin/env bash
set -uo pipefail

ANSWERS=/workspace/answers.env
PIPELINE=/workspace/pipeline.yml
fail() { echo "FAIL: $1"; exit 1; }

[ -f "$PIPELINE" ] || fail "/workspace/pipeline.yml is missing — fixture was tampered with"
[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"

# shellcheck disable=SC1090
source "$ANSWERS"

trim() { echo "$1" | sed -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//'; }
sort_csv() { echo "$1" | tr ',' '\n' | sed -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//' | sed '/^$/d' | sort | tr '\n' ',' | sed 's/,$//'; }

# --- derive ground truth from the fixture itself, never hardcoded ---

# Triggers: 2-space-indented keys under "on:", before "jobs:" starts.
expected_triggers=$(awk '
  /^on:[[:space:]]*$/ { f=1; next }
  /^jobs:[[:space:]]*$/ { f=0 }
  f && /^  [A-Za-z0-9_-]+:/ { line=$0; sub(/^  /,"",line); sub(/:.*/,"",line); print line }
' "$PIPELINE")
[ -n "$expected_triggers" ] || fail "couldn't find any triggers under on: in pipeline.yml — fixture was tampered with"
expected_triggers_csv=$(sort_csv "$(echo "$expected_triggers" | tr '\n' ',')")

# Job names and their needs: value, one "job|needs" pair per line.
job_needs=$(awk '
  /^jobs:[[:space:]]*$/ { f=1; next }
  f && /^  [A-Za-z0-9_-]+:[[:space:]]*$/ {
    job=$0; sub(/^  /,"",job); sub(/:[[:space:]]*$/,"",job); cur=job; needs[cur]=""; order[++n]=cur; next
  }
  f && cur != "" && /needs:/ {
    val=$0
    sub(/^[[:space:]]*needs:[[:space:]]*/,"",val)
    needs[cur]=val
  }
  END { for (i=1;i<=n;i++) print order[i] "|" needs[order[i]] }
' "$PIPELINE")
[ -n "$job_needs" ] || fail "couldn't find any jobs under jobs: in pipeline.yml — fixture was tampered with"

parallel_jobs=()
waiting_job=""
waiting_on=""
while IFS='|' read -r job val; do
  if [ -z "$val" ]; then
    parallel_jobs+=("$job")
  else
    clean=$(echo "$val" | tr -d '[]' | tr ',' '\n' | sed -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//' | sed '/^$/d')
    count=$(echo "$clean" | grep -c .)
    if [ "$count" -eq 1 ]; then
      waiting_job="$job"
      waiting_on="$clean"
    fi
  fi
done <<< "$job_needs"

[ "${#parallel_jobs[@]}" -ge 2 ] || fail "expected at least two jobs with no needs: in pipeline.yml — fixture was tampered with"
[ -n "$waiting_job" ] || fail "expected exactly one job whose needs: names exactly one other job — fixture was tampered with"

expected_parallel_csv=$(sort_csv "$(IFS=,; echo "${parallel_jobs[*]}")")

# Cache key file: the filename passed to hashFiles(...).
expected_cache_file=$(grep -o "hashFiles('[^']*')" "$PIPELINE" | head -1 | sed -e "s/hashFiles('//" -e "s/')$//")
[ -n "$expected_cache_file" ] || fail "couldn't find a hashFiles(...) cache key in pipeline.yml — fixture was tampered with"

# Job with continue-on-error: true.
expected_masked_job=$(awk '
  /^jobs:[[:space:]]*$/ { f=1; next }
  f && /^  [A-Za-z0-9_-]+:[[:space:]]*$/ {
    job=$0; sub(/^  /,"",job); sub(/:[[:space:]]*$/,"",job); cur=job
  }
  f && cur != "" && /continue-on-error:[[:space:]]*true/ { print cur }
' "$PIPELINE" | head -1)
[ -n "$expected_masked_job" ] || fail "couldn't find a job with continue-on-error: true in pipeline.yml — fixture was tampered with"

# Step whose run: command ends in || true — report the preceding name:.
expected_swallowed_step=$(awk '
  /- name:/ { name=$0; sub(/^[[:space:]]*-[[:space:]]*name:[[:space:]]*/,"",name) }
  /\|\| true/ { print name }
' "$PIPELINE" | head -1)
[ -n "$expected_swallowed_step" ] || fail "couldn't find any step ending in || true in pipeline.yml — fixture was tampered with"

# --- grade the learner's answers against that ground truth ---

[ -n "${TRIGGER_EVENTS:-}" ] || fail "TRIGGER_EVENTS is empty"
got_triggers_csv=$(sort_csv "$TRIGGER_EVENTS")
[ "$got_triggers_csv" = "$expected_triggers_csv" ] || fail "TRIGGER_EVENTS should list: $expected_triggers_csv (got: $got_triggers_csv)"

[ -n "${PARALLEL_JOBS:-}" ] || fail "PARALLEL_JOBS is empty"
got_parallel_csv=$(sort_csv "$PARALLEL_JOBS")
[ "$got_parallel_csv" = "$expected_parallel_csv" ] || fail "PARALLEL_JOBS should list: $expected_parallel_csv (got: $got_parallel_csv)"

[ -n "${WAITING_JOB:-}" ] || fail "WAITING_JOB is empty"
[ "$(trim "$WAITING_JOB")" = "$waiting_job" ] || fail "WAITING_JOB should be $waiting_job (got: $WAITING_JOB)"

[ -n "${WAITING_ON:-}" ] || fail "WAITING_ON is empty"
[ "$(trim "$WAITING_ON")" = "$waiting_on" ] || fail "WAITING_ON should be $waiting_on (got: $WAITING_ON)"

[ -n "${CACHE_KEY_FILE:-}" ] || fail "CACHE_KEY_FILE is empty"
[ "$(trim "$CACHE_KEY_FILE")" = "$expected_cache_file" ] || fail "CACHE_KEY_FILE should be $expected_cache_file (got: $CACHE_KEY_FILE)"

[ -n "${MASKED_FAILURE_JOB:-}" ] || fail "MASKED_FAILURE_JOB is empty"
[ "$(trim "$MASKED_FAILURE_JOB")" = "$expected_masked_job" ] || fail "MASKED_FAILURE_JOB should be $expected_masked_job (got: $MASKED_FAILURE_JOB)"

[ -n "${SWALLOWED_STEP:-}" ] || fail "SWALLOWED_STEP is empty"
[ "$(trim "$SWALLOWED_STEP")" = "$expected_swallowed_step" ] || fail "SWALLOWED_STEP should be $expected_swallowed_step (got: $SWALLOWED_STEP)"

echo "PASS: answers.env correctly reads the pipeline's triggers, job ordering, cache key, and swallowed failures"
