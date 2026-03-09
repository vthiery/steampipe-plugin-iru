#!/usr/bin/env bash
# Run a smoke-test query against every iru table and report results.
# Exit code is the number of failed tables (0 = all passed).

set -euo pipefail

# Colour codes (disabled if not a terminal)
if [ -t 1 ]; then
  GREEN=$'\033[0;32m'
  RED=$'\033[0;31m'
  YELLOW=$'\033[0;33m'
  RESET=$'\033[0m'
else
  GREEN="" RED="" YELLOW="" RESET=""
fi

PASS=0
FAIL=0
SKIP=0

run_test() {
  local table="$1"
  local query="$2"
  printf "  %-40s" "$table"
  local output
  if output=$(steampipe query "$query" 2>&1); then
    printf "${GREEN}PASS${RESET}\n"
    ((PASS++)) || true
  else
    # A 403 / authorization error is treated as a skip, not a hard failure
    if echo "$output" | grep -qiE "403|unauthorized|forbidden"; then
      printf "${YELLOW}SKIP${RESET} (authorization error)\n"
      ((SKIP++)) || true
    else
      printf "${RED}FAIL${RESET}\n"
      echo "$output" | sed 's/^/    /'
      ((FAIL++)) || true
    fi
  fi
}

echo ""
echo "Iru Steampipe plugin — table smoke tests"
echo "========================================="
echo ""

run_test "iru_device"       "select device_id, device_name, serial_number from iru_device limit 1"
run_test "iru_user"         "select id, name, email from iru_user limit 1"
run_test "iru_blueprint"    "select id, name, active_devices from iru_blueprint limit 1"
run_test "iru_tag"          "select id, name, color from iru_tag limit 1"
run_test "iru_threat"       "select threat_id, device_name, classification from iru_threat limit 1"
run_test "iru_library_item" "select id, name, item_type from iru_library_item limit 1"
run_test "iru_audit_log"    "select id, event_type, actor_name from iru_audit_log limit 1"

echo ""
echo "Results: ${GREEN}${PASS} passed${RESET}  ${YELLOW}${SKIP} skipped${RESET}  ${RED}${FAIL} failed${RESET}"
echo ""

exit "$FAIL"
