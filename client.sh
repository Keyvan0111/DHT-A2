#!/usr/bin/env bash
set -euo pipefail

# Simple HTTP client for the DHT using curl
# Usage:
#   ./client.sh --entry host:port --phase put --num-ops 1000 --keys-seed 42 --value-size 128
#
# Phases:
#   put -> PUT /storage/<key> with <value>
#   get -> GET /storage/<key>
#
# Notes:
# - Keys are deterministic from --keys-seed to allow a GET run after a PUT run.
# - Value is a fixed-size ASCII blob (default 64 bytes).

ENTRY=""
PHASE="put"
NUM_OPS=100
KEYS_SEED=0
VALUE_SIZE=64

while [[ $# -gt 0 ]]; do
  case "$1" in
    --entry) ENTRY="$2"; shift 2;;
    --phase) PHASE="$2"; shift 2;;
    --num-ops) NUM_OPS="$2"; shift 2;;
    --keys-seed) KEYS_SEED="$2"; shift 2;;
    --value-size) VALUE_SIZE="$2"; shift 2;;
    *) echo "Unknown arg: $1" >&2; exit 2;;
  esac
done

if [[ -z "$ENTRY" ]]; then
  echo "Missing --entry host:port" >&2
  exit 2
fi

# Fixed-size ASCII value
VAL="$(head -c "$VALUE_SIZE" < /dev/zero | tr '\0' 'A')"

# Generate deterministic keys (k<seed>_<i>)
start="$KEYS_SEED"
end=$(( KEYS_SEED + NUM_OPS - 1 ))

case "$PHASE" in
  put)
    for i in $(seq "$start" "$end"); do
      key="k${i}"
      # PUT value
      # Important: send raw body (no JSON), Content-Type text/plain
      curl -sS -X PUT "http://$ENTRY/storage/$key" \
           -H 'Content-Type: text/plain' \
           --data-binary "$VAL" > /dev/null
    done
    ;;
  get)
    for i in $(seq "$start" "$end"); do
      key="k${i}"
      curl -sS "http://$ENTRY/storage/$key" > /dev/null || true
    done
    ;;
  *)
    echo "Unknown --phase $PHASE (expected put|get)" >&2
    exit 2
    ;;
esac
