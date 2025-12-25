#!/usr/bin/env bash
set -euo pipefail

base_url=${1:-http://127.0.0.1:4010}

for _ in {1..20}; do
  if curl -sf "$base_url/prices/singles" >/dev/null; then
    break
  fi
  sleep 1
done

curl -sf "$base_url/prices/singles" >/dev/null
curl -sf "$base_url/prices/variants" >/dev/null
curl -sf "$base_url/prices/sealed" >/dev/null

curl -sf \
  -H "X-ManaPool-Access-Token: test-token" \
  -H "X-ManaPool-Email: test@example.com" \
  "$base_url/account" >/dev/null

curl -sf \
  -H "Content-Type: application/json" \
  -d '{"cart":[{"type":"mtg_single","name":"Polar Kraken","quantity_requested":1,"language_ids":["EN"],"finish_ids":["NF"],"condition_ids":["NM"]}]}' \
  "$base_url/buyer/optimizer" >/dev/null

curl -sf \
  -H "Content-Type: application/json" \
  -d '{"commander_names":["Atraxa, Praetors\u0027 Voice"],"other_cards":[{"name":"Lightning Bolt","quantity":4}]}' \
  "$base_url/deck" >/dev/null

curl -sf \
  -H "Content-Type: application/json" \
  -d '{"card_names":["Lightning Bolt"]}' \
  "$base_url/card_info" >/dev/null

curl -sf \
  -H "X-ManaPool-Access-Token: test-token" \
  -H "X-ManaPool-Email: test@example.com" \
  "$base_url/buyer/orders?since=2024-04-01T00:00:00Z&limit=1" >/dev/null
