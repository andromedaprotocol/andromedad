#!/bin/sh
#set -o errexit -o nounset -o pipefail

PASSWORD=${PASSWORD:-1234567890}
ANDR=${STAKE_TOKEN:-uandr}
CHAIN_ID=${CHAIN_ID:-testing}
MONIKER=${MONIKER:-node001}

andromedad init --chain-id "$CHAIN_ID" "$MONIKER"
# staking/governance token is hardcoded in config, change this
sed -i "s/\"uandr\"/\"$ANDR\"/" "$HOME"/.andromedad/config/genesis.json
# this is essential for sub-1s block times (or header times go crazy)
sed -i 's/"time_iota_ms": "1000"/"time_iota_ms": "10"/' "$HOME"/.andromedad/config/genesis.json

if ! andromedad keys show validator; then
  (echo "$PASSWORD"; echo "$PASSWORD") | andromedad keys add validator
fi
# hardcode the validator account for this instance
echo "$PASSWORD" | andromedad add-genesis-account validator "1000000000$ANDR,1000000000$ANDR"

# (optionally) add a few more genesis accounts
for addr in "$@"; do
  echo $addr
  andromedad add-genesis-account "$addr" "1000000000$ANDR,1000000000$ANDR"
done

# submit a genesis validator tx
## Workraround for https://github.com/cosmos/cosmos-sdk/issues/8251
(echo "$PASSWORD"; echo "$PASSWORD"; echo "$PASSWORD") | andromedad gentx validator "250000000$ANDR" --chain-id="$CHAIN_ID" --amount="250000000$ANDR"
## should be:
# (echo "$PASSWORD"; echo "$PASSWORD"; echo "$PASSWORD") | andromedad gentx validator "250000000$ANDR" --chain-id="$CHAIN_ID"
andromedad collect-gentxs
