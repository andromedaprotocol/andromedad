#!/bin/bash
set -o errexit -o nounset -o pipefail

BASE_ACCOUNT=$(andromedad keys show validator -a)
andromedad q account "$BASE_ACCOUNT" -o json | jq

echo "## Add new account"
andromedad keys add fred

echo "## Check balance"
NEW_ACCOUNT=$(andromedad keys show fred -a)
andromedad q bank balances "$NEW_ACCOUNT" -o json || true

echo "## Transfer tokens"
andromedad tx bank send validator "$NEW_ACCOUNT" 1ustake --gas 1000000 -y --chain-id=testing --node=http://localhost:26657 -b block -o json | jq

echo "## Check balance again"
andromedad q bank balances "$NEW_ACCOUNT" -o json | jq
