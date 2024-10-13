package helpers

import (
	"context"
	"strings"
	"testing"

	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/stretchr/testify/require"
)

type (
	GetCountQuery struct {
		// {"get_total_funds":{"addr":"andr1..."}}
		Addr string `json:"addr"`
	}
	GetTotalFundsQuery struct {
		// {"get_total_funds":{"addr":"andr1..."}}
		Addr string `json:"addr"`
	}
	GetTotalFundsResponse struct {
		// {"data":{"total_funds":[{"denom":"ibc/04F5F501207C3626A2C14BFEF654D51C2E0B8F7CA578AB8ED272A66FE4E48097","amount":"1"}]}}
		Data *GetTotalFundsObj `json:"data"`
	}
	GetTotalFundsObj struct {
		TotalFunds []WasmCoin `json:"total_funds"`
	}
	GetCountResponse struct {
		// {"data":{"count":0}}
		Data *GetCountObj `json:"data"`
	}
	GetCountObj struct {
		Count int64 `json:"count"`
	}
)

func GetIBCHooksUserAddress(t *testing.T, ctx context.Context, chain *cosmos.CosmosChain, channel, uaddr string) string {
	// appd q ibchooks wasm-sender channel-0 "andr11hj5fveer5cjtn4wd6wstzugjfdxzl0xps73ftl" --node http://localhost:26657
	cmd := []string{
		chain.Config().Bin, "query", "ibchooks", "wasm-sender", channel, uaddr,
		"--node", chain.GetRPCAddress(),
		"--output", "json",
	}

	// This query does not return a type, just prints the string.
	stdout, _, err := chain.Exec(ctx, cmd, nil)
	require.NoError(t, err)

	address := strings.Replace(string(stdout), "\n", "", -1)
	return address
}

func GetIBCHookTotalFunds(t *testing.T, ctx context.Context, chain *cosmos.CosmosChain, contract string, uaddr string) GetTotalFundsResponse {
	var res GetTotalFundsResponse
	err := chain.QueryContract(ctx, contract, QueryMsg{GetTotalFunds: &GetTotalFundsQuery{Addr: uaddr}}, &res)
	require.NoError(t, err)
	return res
}

func GetIBCHookCount(t *testing.T, ctx context.Context, chain *cosmos.CosmosChain, contract string, uaddr string) GetCountResponse {
	var res GetCountResponse
	err := chain.QueryContract(ctx, contract, QueryMsg{GetCount: &GetCountQuery{Addr: uaddr}}, &res)
	require.NoError(t, err)
	return res
}
