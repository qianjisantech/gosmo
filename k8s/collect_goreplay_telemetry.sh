#!/usr/bin/env bash
#
# collect_goreplay_telemetry.sh
#
# Gathers telemetry from a GoReplay DaemonSet in the 'goreplay' namespace.
# Works on macOS and Linux, assuming 'kubectl' (or compatible) is installed.
#
# Usage examples:
#   ./collect_goreplay_telemetry.sh
#   ./collect_goreplay_telemetry.sh "microk8s kubectl"

set -euo pipefail

########################################
# Determine kubectl command
########################################
if [[ $# -gt 0 ]]; then
  # If an argument was provided, use that as the kubectl command
  KUBECTL="$*"
else
  # Default to 'kubectl'
  KUBECTL="kubectl"
fi

########################################
# Check that the base command exists
########################################
# For "microk8s kubectl", we only check "microk8s" in PATH. For "oc" we check "oc".
BASE_CMD="${KUBECTL%% *}"  # everything before the first space
if ! command -v "${BASE_CMD}" >/dev/null 2>&1; then
  echo "ERROR: '${BASE_CMD}' not found in PATH. Please install or configure it first."
  exit 1
fi

echo "Using kubectl command: $KUBECTL"
echo

########################################
# Helper function to print and run commands
########################################
run_cmd() {
  echo "Command: $*"
  eval "$*"
}

########################################
# 1. Print logs from ALL GoReplay pods
########################################
echo "=================================================="
echo "1. Gathering logs from all goreplay pods (all containers)..."
echo "=================================================="
run_cmd "$KUBECTL logs -n goreplay -l app=goreplay --all-containers" || {
  echo "WARNING: Failed to get logs from pods with label app=goreplay"
}

########################################
# 2. Describe the GoReplay DaemonSet
########################################
echo
echo "=================================================="
echo "2. Describing DaemonSet goreplay-daemon..."
echo "=================================================="
run_cmd "$KUBECTL describe daemonset goreplay-daemon -n goreplay" || {
  echo "WARNING: Failed to describe daemonset goreplay-daemon"
}

########################################
# 3. Get list of GoReplay pods (full output)
########################################
echo
echo "=================================================="
echo "3. Listing goreplay pods (full output)..."
echo "=================================================="

# Print full output (no -o name here):
run_cmd "$KUBECTL get pods -n goreplay -l app=goreplay"

# Then retrieve just the names for further processing:
echo
echo "Getting goreplay pod names for telemetry collection..."
pods=$($KUBECTL get pods -n goreplay -l app=goreplay -o name 2>/dev/null) || {
  echo "ERROR: Failed to list pods with label app=goreplay"
  exit 1
}
echo "Found pods:"
echo "$pods"
echo

########################################
# 4. For each pod, gather logs, describe, and get events
########################################
for pod in $pods; do
  # pod looks like "pod/goreplay-daemon-xyz"
  pod_name="${pod##*/}"  # remove "pod/" prefix

  echo "=================================================="
  echo "LOGS for pod: ${pod_name}"
  echo "=================================================="
  run_cmd "$KUBECTL logs ${pod_name} -n goreplay" || {
    echo "WARNING: Failed to get logs for pod ${pod_name}"
  }

  echo
  echo "--------------------------------------------------"
  echo "DESCRIBE for pod: ${pod_name}"
  echo "--------------------------------------------------"
  run_cmd "$KUBECTL describe pod -n goreplay ${pod_name}" || {
    echo "WARNING: Failed to describe pod ${pod_name}"
  }

  echo
  echo "--------------------------------------------------"
  echo "EVENTS for pod: ${pod_name}"
  echo "--------------------------------------------------"
  run_cmd "$KUBECTL get events -n goreplay --field-selector involvedObject.name=${pod_name}" || {
    echo "WARNING: Failed to get events for pod ${pod_name}"
  }
  echo
done

echo "=================================================="
echo "Telemetry collection complete."
echo "=================================================="
