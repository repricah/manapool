#!/usr/bin/env bash
set -euo pipefail

base_url=${1:-http://127.0.0.1:4010}
summary_file="e2e-summary.md"
failed_tests=0
total_tests=0

# Initialize summary
echo "# E2E Mock Test Results" > "$summary_file"
echo "" >> "$summary_file"

test_get() {
  local endpoint=$1
  local test_name="GET $endpoint"
  echo "==> Testing $test_name"
  
  total_tests=$((total_tests + 1))
  
  local status_code="000"
  set +e  # Temporarily disable exit on error
  status_code=$(curl -sS -w "%{http_code}" -o /dev/null "$base_url$endpoint" 2>/dev/null)
  set -e  # Re-enable exit on error
  
  if [[ "$status_code" =~ ^2[0-9][0-9]$ ]]; then
    echo "✅ $test_name - Status: $status_code" >> "$summary_file"
    echo "    ✅ Status: $status_code"
  else
    echo "❌ $test_name - Status: $status_code" >> "$summary_file"
    echo "    ❌ Status: $status_code"
    failed_tests=$((failed_tests + 1))
  fi
}

test_post() {
  local endpoint=$1
  local data=$2
  local test_name="POST $endpoint"
  echo "==> Testing $test_name"
  
  total_tests=$((total_tests + 1))
  
  local status_code="000"
  set +e  # Temporarily disable exit on error
  status_code=$(curl -sS -w "%{http_code}" -o /dev/null \
    -H "Content-Type: application/json" \
    -d "$data" \
    -X POST "$base_url$endpoint" 2>/dev/null)
  set -e  # Re-enable exit on error
  
  if [[ "$status_code" =~ ^2[0-9][0-9]$ ]]; then
    echo "✅ $test_name - Status: $status_code" >> "$summary_file"
    echo "    ✅ Status: $status_code"
  else
    echo "❌ $test_name - Status: $status_code" >> "$summary_file"
    echo "    ❌ Status: $status_code"
    failed_tests=$((failed_tests + 1))
  fi
}

# Wait for Prism to be ready
echo "Waiting for Prism to be ready..."
ready=false
for i in {1..30}; do
  if curl -sf "$base_url/prices/singles" >/dev/null 2>&1; then
    ready=true
    echo "✅ Prism is ready"
    break
  fi
  sleep 2
done

if [[ "$ready" != "true" ]]; then
  echo "❌ Prism mock server failed to start"
  exit 1
fi

# Run tests
test_get "/prices/singles"
test_get "/prices/variants"  
test_get "/prices/sealed"
test_post "/buyer/optimizer" '{"cart":[{"type":"mtg_single","name":"Polar Kraken","quantity_requested":1,"language_ids":["EN"],"finish_ids":["NF"],"condition_ids":["NM"]}]}'

# Generate final summary
echo "" >> "$summary_file"
passed_tests=$((total_tests - failed_tests))

echo "## Summary" >> "$summary_file"
echo "- Total tests: $total_tests" >> "$summary_file"
echo "- Passed: $passed_tests" >> "$summary_file"
echo "- Failed: $failed_tests" >> "$summary_file"

echo ""
echo "Summary: $passed_tests/$total_tests tests passed"

if [[ "$failed_tests" -gt 0 ]]; then
  echo "❌ $failed_tests tests failed"
  exit 1
fi

echo "✅ All tests passed!"
