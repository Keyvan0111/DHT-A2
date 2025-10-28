#!/usr/bin/env bash
set -euo pipefail

# Run throughput experiments across cluster sizes and write CSV.
#
# Usage:
#   ./run_bench.sh results.csv
#
# Optional env vars:
#   K=1000           # operations per phase (PUT/GET)
#   RUNS=3           # repetitions per N
#   SIZES="1 2 4 8 16"  # cluster sizes to test
#   VALUE_SIZE=64    # bytes per value for PUT
#   SETTLE=3         # seconds to wait after launch
#
# Requirements:
# - Assumes your ./run.sh <N> prints a JSON array of "host:port" entries.
# - Uses the first entry as the client gateway (any node should forward).
# - Uses ./client.sh to drive PUT/GET phases.
#
# CSV columns:
#   N,PHASE,run_idx,num_ops,seconds,throughput_ops_per_sec
#
# Notes:
# - You may need to manually clean up processes on the cluster after runs.
# - For fairness, we do PUT first (to populate), then GET over same keys.

OUT="${1:-results.csv}"

K="${K:-1000}"
RUNS="${RUNS:-3}"
SIZES="${SIZES:-1 2 4 8 16}"
VALUE_SIZE="${VALUE_SIZE:-64}"
SETTLE="${SETTLE:-3}"

# Ensure client exists
if [[ ! -x "./client.sh" ]]; then
  echo "client.sh not found or not executable. Place it next to this script." >&2
  exit 3
fi

# Header if file doesn't exist
if [[ ! -f "$OUT" ]]; then
  echo "N,PHASE,run_idx,num_ops,seconds,throughput_ops_per_sec" > "$OUT"
fi

for N in $SIZES; do
  echo "== Launching $N nodes ==" >&2
  HOSTS_JSON="$(./run.sh "$N")"   # e.g., ["node1:50001","node2:50033",...]
  ENTRY="$(echo "$HOSTS_JSON" | sed -E 's/.*"\s*([^"]+)".*/\1/')"  # pick first "host:port"
  if [[ -z "$ENTRY" ]]; then
    echo "Failed to parse entry from run.sh output: $HOSTS_JSON" >&2
    exit 4
  fi
  echo "Entry node: $ENTRY" >&2
  sleep "$SETTLE"

  for r in $(seq 1 "$RUNS"); do
    echo "-- N=$N, run $r: PUT K=$K --" >&2
    t_put_start=$(date +%s.%N)
    ./client.sh --entry "$ENTRY" --phase put --num-ops "$K" --keys-seed 42 --value-size "$VALUE_SIZE"
    t_put_end=$(date +%s.%N)
    put_secs=$(echo "$t_put_end - $t_put_start" | bc -l)
    put_thr=$(echo "$K / $put_secs" | bc -l)
    printf "%s,%s,%s,%s,%.6f,%.6f\n" "$N" "PUT" "$r" "$K" "$put_secs" "$put_thr" >> "$OUT"

    echo "-- N=$N, run $r: GET K=$K --" >&2
    t_get_start=$(date +%s.%N)
    ./client.sh --entry "$ENTRY" --phase get --num-ops "$K" --keys-seed 42
    t_get_end=$(date +%s.%N)
    get_secs=$(echo "$t_get_end - $t_get_start" | bc -l)
    get_thr=$(echo "$K / $get_secs" | bc -l)
    printf "%s,%s,%s,%s,%.6f,%.6f\n" "$N" "GET" "$r" "$K" "$get_secs" "$get_thr" >> "$OUT"
  done

  echo "== Completed N=$N ==" >&2
done

echo "Done. Wrote $OUT" >&2
