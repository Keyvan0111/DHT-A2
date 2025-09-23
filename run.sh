#!/usr/bin/env bash
set -euo pipefail

NUM="${1}"

# ---------- config ----------
SERVER_BIN="./main"                                 # compiled Go binary path
LOG_DIR="./logs"                                    # logs + port files
AVAILABLE_NODES_CMD="/share/ifi/available-nodes.sh" # provided by UiT
PORT_WAIT_SECS="10"                                 # per-server wait

PROJECT_DIR="$(cd "$(dirname "$0")" && pwd)"
mkdir -p "$PROJECT_DIR/$LOG_DIR"

# Build the Go server
( cd "$PROJECT_DIR" && go build -o "${SERVER_BIN#./}" ./cmd/app/server.go )

# ---------- get nodes ----------
declare -a NODES=()
if [[ -x "$AVAILABLE_NODES_CMD" ]]; then
  # Deduplicate, skip empties
  while IFS= read -r line; do
    [[ -n "$line" ]] && NODES+=("$line")
  done < <("$AVAILABLE_NODES_CMD" | awk 'NF && !seen[$0]++')
fi

if ((${#NODES[@]} == 0)); then
  echo "No nodes available." >&2
  exit 4
fi

for ((i = 0; i<NUM; i++)); do
  node="${NODES[$(( i % ${#NODES[@]} ))]}"
  
  port=""
    port="$(cat "$port_file" 2>/dev/null || true)"

# ---------- launch & collect ----------
declare -a HOSTPORTS=()

for ((i=0; i<NUM; i++)); do
  node="${NODES[$(( i % ${#NODES[@]} ))]}"

  port_file="$PROJECT_DIR/$LOG_DIR/port_${node}_$$_${i}.txt"
  log_file="$PROJECT_DIR/$LOG_DIR/server_${node}_$$_${i}.log"

  # Launch detached on the compute node; stdout/err go to per-server log file on the shared NFS
  ssh -o BatchMode=yes -o StrictHostKeyChecking=no -f "$node" \
    "cd '$PROJECT_DIR' && PORT_FILE='$port_file' nohup '$SERVER_BIN' > '$log_file' 2>&1 < /dev/null &"

  # Poll for port file
  # We did this to be able to check the logs of our server on each clusternode to validate
  port=""
  for ((t=0; t<PORT_WAIT_SECS*10; t++)); do
    if [[ -s "$port_file" ]]; then
      port="$(cat "$port_file" 2>/dev/null || true)"
      if [[ "$port" =~ ^[0-9]+$ ]]; then
        break
      fi
    fi
    sleep 0.1
  done

  if [[ -z "$port" || ! "$port" =~ ^[0-9]+$ ]]; then
    echo "Failed to obtain port from $port_file (node $node). Check $log_file" >&2
    exit 5
  fi

  HOSTPORTS+=("${node}:${port}")
done

# ---------- print JSON list ----------
printf '['
for ((i=0; i<${#HOSTPORTS[@]}; i++)); do
  (( i > 0 )) && printf ','
  printf '"%s"' "${HOSTPORTS[$i]}"
done
printf ']\n'
