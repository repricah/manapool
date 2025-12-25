#!/usr/bin/env bash
set -euo pipefail

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
    "${url}"
  echo "    body: $(cat /tmp/e2e-body.json)"
}

ready=false
for _ in {1..20}; do
  if curl -sf "$base_url/prices/singles" >/dev/null; then
    ready=true
    break
  fi
  sleep 1
done

if [[ "${ready}" != "true" ]]; then
  echo "Prism mock did not become ready at ${base_url}"
  exit 1
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
