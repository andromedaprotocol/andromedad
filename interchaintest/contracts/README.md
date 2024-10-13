# Contracts

A list of the contracts here which are pre-compiled in other repos.

> cw_template -> `cargo generate --git https://github.com/CosmWasm/cw-template.git --name cwtemplate -d minimal=false --tag 9ef93954d62f383468bb1ec869fef894be53bc4d`

cw_template uses an older version so cosmwasm_1_3 feature is not required (breaks upgrade test since old v1.1.0-beta1 version did not have it. Same concept though)

> ibchooks_counter.wasm -> <https://github.com/osmosis-labs/osmosis/blob/64393a14e18b2562d72a3892eec716197a3716c7/tests/ibc-hooks/bytecode/counter.wasm>

SubMessage test
> cw721_base      - https://github.com/CosmWasm/cw-nfts/releases/download/v0.17.0/cw721_base.wasm

> cw721-receiver  - https://github.com/CosmWasm/cw-nfts/pull/144
