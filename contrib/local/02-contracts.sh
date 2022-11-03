#!/bin/bash
set -o errexit -o nounset -o pipefail -x

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"

echo "-----------------------"
echo "## Add new CosmWasm contract"
RESP=$(andromedad tx wasm store "$DIR/../../x/wasm/keeper/testdata/hackatom.wasm" \
  --from validator --gas 1500000 -y --chain-id=testing --node=http://localhost:26657 -b block -o json)

CODE_ID=$(echo "$RESP" | jq -r '.logs[0].events[1].attributes[-1].value')
CODE_HASH=$(echo "$RESP" | jq -r '.logs[0].events[1].attributes[-2].value')
echo "* Code id: $CODE_ID"
echo "* Code checksum: $CODE_HASH"

echo "* Download code"
TMPDIR=$(mktemp -t andromedadXXXXXX)
andromedad q wasm code "$CODE_ID" "$TMPDIR"
rm -f "$TMPDIR"
echo "-----------------------"
echo "## List code"
andromedad query wasm list-code --node=http://localhost:26657 --chain-id=testing -o json | jq

echo "-----------------------"
echo "## Create new contract instance"
INIT="{\"verifier\":\"$(andromedad keys show validator -a)\", \"beneficiary\":\"$(andromedad keys show fred -a)\"}"
andromedad tx wasm instantiate "$CODE_ID" "$INIT" --admin="$(andromedad keys show validator -a)" \
  --from validator --amount="100ustake" --label "local0.1.0" \
  --gas 1000000 -y --chain-id=testing -b block -o json | jq

CONTRACT=$(andromedad query wasm list-contract-by-code "$CODE_ID" -o json | jq -r '.contracts[-1]')
echo "* Contract address: $CONTRACT"

echo "## Create new contract instance with predictable address"
andromedad tx wasm instantiate2 "$CODE_ID" "$INIT" $(echo -n "testing" | xxd -ps) \
  --admin="$(andromedad keys show validator -a)" \
  --from validator --amount="100ustake" --label "local0.1.0" \
  --fix-msg \
  --gas 1000000 -y --chain-id=testing -b block -o json | jq

predictedAdress=$(andromedad q wasm build-address "$CODE_HASH" $(andromedad keys show validator -a) $(echo -n "testing" | xxd -ps) "$INIT")
andromedad q wasm contract "$predictedAdress" -o json | jq

echo "### Query all"
RESP=$(andromedad query wasm contract-state all "$CONTRACT" -o json)
echo "$RESP" | jq
echo "### Query smart"
andromedad query wasm contract-state smart "$CONTRACT" '{"verifier":{}}' -o json | jq
echo "### Query raw"
KEY=$(echo "$RESP" | jq -r ".models[0].key")
andromedad query wasm contract-state raw "$CONTRACT" "$KEY" -o json | jq

echo "-----------------------"
echo "## Execute contract $CONTRACT"
MSG='{"release":{}}'
andromedad tx wasm execute "$CONTRACT" "$MSG" \
  --from validator \
  --gas 1000000 -y --chain-id=testing -b block -o json | jq

echo "-----------------------"
echo "## Set new admin"
echo "### Query old admin: $(andromedad q wasm contract "$CONTRACT" -o json | jq -r '.contract_info.admin')"
echo "### Update contract"
andromedad tx wasm set-contract-admin "$CONTRACT" "$(andromedad keys show fred -a)" \
  --from validator -y --chain-id=testing -b block -o json | jq
echo "### Query new admin: $(andromedad q wasm contract "$CONTRACT" -o json | jq -r '.contract_info.admin')"

echo "-----------------------"
echo "## Migrate contract"
echo "### Upload new code"
RESP=$(andromedad tx wasm store "$DIR/../../x/wasm/keeper/testdata/burner.wasm" \
  --from validator --gas 1000000 -y --chain-id=testing --node=http://localhost:26657 -b block -o json)

BURNER_CODE_ID=$(echo "$RESP" | jq -r '.logs[0].events[1].attributes[-1].value')
echo "### Migrate to code id: $BURNER_CODE_ID"

DEST_ACCOUNT=$(andromedad keys show fred -a)
andromedad tx wasm migrate "$CONTRACT" "$BURNER_CODE_ID" "{\"payout\": \"$DEST_ACCOUNT\"}" --from fred \
  --chain-id=testing -b block -y -o json | jq

echo "### Query destination account: $BURNER_CODE_ID"
andromedad q bank balances "$DEST_ACCOUNT" -o json | jq
echo "### Query contract meta data: $CONTRACT"
andromedad q wasm contract "$CONTRACT" -o json | jq

echo "### Query contract meta history: $CONTRACT"
andromedad q wasm contract-history "$CONTRACT" -o json | jq

echo "-----------------------"
echo "## Clear contract admin"
echo "### Query old admin: $(andromedad q wasm contract "$CONTRACT" -o json | jq -r '.contract_info.admin')"
echo "### Update contract"
andromedad tx wasm clear-contract-admin "$CONTRACT" \
  --from fred -y --chain-id=testing -b block -o json | jq
echo "### Query new admin: $(andromedad q wasm contract "$CONTRACT" -o json | jq -r '.contract_info.admin')"
