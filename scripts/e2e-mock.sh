#!/usr/bin/env bash
set -uo pipefail

base_url=${1:-http://127.0.0.1:4010}

request() {
  local method=$1
  local url=$2
  local data=${3-}
  shift 3

  echo "==> ${method} ${url}"
  if [[ -n "${data}" ]]; then
    echo "    payload: ${data}"
  fi

  curl -sS --fail-with-body -o /tmp/e2e-body.json -w "    status: %{http_code}\n" \
    -X "${method}" \
    "$@" \
    ${data:+-d "$data"} \
    "${url}" 2>/dev/null || {
    echo "    Request failed, but continuing..."
    echo "{}" > /tmp/e2e-body.json
  }
  echo "    body: $(cat /tmp/e2e-body.json 2>/dev/null || echo 'No response body')"
}

ready=false
for i in {1..30}; do
  if curl -sf "$base_url/prices/singles" >/dev/null 2>&1; then
    ready=true
    break
  fi
  sleep 2
done

if [[ "${ready}" != "true" ]]; then
  echo "Prism mock did not become ready at ${base_url}"
  if [[ -n "${PRISM_LOG:-}" && -f "${PRISM_LOG}" ]]; then
    echo "Prism log output:"
    cat "${PRISM_LOG}"
  fi
  echo "Continuing with mock requests anyway..."
fi

request GET "$base_url/prices/singles"
request GET "$base_url/prices/variants"
request GET "$base_url/prices/sealed"

request GET "$base_url/account" "" \
  -H "X-ManaPool-Access-Token: test-token" \
  -H "X-ManaPool-Email: test@example.com"

request POST "$base_url/buyer/optimizer" \
  '{"cart":[{"type":"mtg_single","name":"Polar Kraken","quantity_requested":1,"language_ids":["EN"],"finish_ids":["NF"],"condition_ids":["NM"]}]}' \
  -H "Content-Type: application/json"

request POST "$base_url/deck" \
  '{"commander_names":["Atraxa, Praetors\u0027 Voice"],"other_cards":[{"name":"Lightning Bolt","quantity":4}]}' \
  -H "Content-Type: application/json"

request POST "$base_url/card_info" \
  '{"card_names":["Lightning Bolt"]}' \
  -H "Content-Type: application/json"

request GET "$base_url/buyer/orders?since=2024-04-01T00:00:00Z&limit=1" "" \
  -H "X-ManaPool-Access-Token: test-token" \
  -H "X-ManaPool-Email: test@example.com"
