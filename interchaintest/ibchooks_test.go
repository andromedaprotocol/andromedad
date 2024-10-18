package e2e

import (
	"context"
	"fmt"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	helpers "github.com/andromedaprotocol/andromedad/interchaintest/helpers"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	interchaintestrelayer "github.com/strangelove-ventures/interchaintest/v8/relayer"
	"github.com/strangelove-ventures/interchaintest/v8/testreporter"
	"github.com/strangelove-ventures/interchaintest/v8/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

func TestIBCHooks(t *testing.T) {
	const (
		path = "ibc-path"
	)

	t.Parallel()

	ctx := context.Background()
	rep := testreporter.NewNopReporter()
	eRep := rep.RelayerExecReporter(t)
	client, network := interchaintest.DockerSetup(t)

	cs := &DefaultChainSpec

	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		cs,
		&SecondDefaultChainSpec,
	})

	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	chainA, chainB := chains[0].(*cosmos.CosmosChain), chains[1].(*cosmos.CosmosChain)

	// Relayer Factory
	r := interchaintest.NewBuiltinRelayerFactory(
		ibc.CosmosRly,
		zaptest.NewLogger(t, zaptest.Level(zapcore.DebugLevel)),
		interchaintestrelayer.CustomDockerImage(RelayerRepo, RelayerVersion, "100:1000"),
		interchaintestrelayer.StartupFlags("--processor", "events", "--block-history", "200"),
	).Build(t, client, network)

	ic := interchaintest.NewInterchain().
		AddChain(chainA).
		AddChain(chainB).
		AddRelayer(r, "relayer")

	ic = ic.AddLink(interchaintest.InterchainLink{
		Chain1:  chainA,
		Chain2:  chainB,
		Relayer: r,
		Path:    ibcPath,
	})

	// Build interchain
	require.NoError(t, ic.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:         t.Name(),
		Client:           client,
		NetworkID:        network,
		SkipPathCreation: false,
	}))

	// Create some user accounts on both chains
	users := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), GenesisFundsAmount, chainA, chainB)

	// Get our Bech32 encoded user addresses
	user1, user2 := users[0], users[1]

	user1Addr := user1.FormattedAddress()

	channel, err := ibc.GetTransferChannel(ctx, r, eRep, chainA.Config().ChainID, chainB.Config().ChainID)
	require.NoError(t, err)

	err = r.StartRelayer(ctx, eRep, path)
	require.NoError(t, err)

	t.Cleanup(
		func() {
			err := r.StopRelayer(ctx, eRep)
			if err != nil {
				t.Logf("an error occurred while stopping the relayer: %s", err)
			}
		},
	)

	_, contractAddr := helpers.SetupContract(t, ctx, chainB, user2.KeyName(), "contracts/ibchooks_counter.wasm", `{"count":0}`)

	// do an ibc transfer through the memo to the other chain.
	transfer := ibc.WalletAmount{
		Address: contractAddr,
		Denom:   chainA.Config().Denom,
		Amount:  sdkmath.NewInt(1),
	}

	memo := ibc.TransferOptions{
		Memo: fmt.Sprintf(`{"wasm":{"contract":"%s","msg":%s}}`, contractAddr, `{"increment":{}}`),
	}

	// Initial transfer. Account is created by the wasm execute is not so we must do this twice to properly set up
	transferTx, err := chainA.SendIBCTransfer(ctx, channel.ChannelID, user1.KeyName(), transfer, memo)
	require.NoError(t, err)
	andrHeight, err := chainA.Height(ctx)
	require.NoError(t, err)

	_, err = testutil.PollForAck(ctx, chainA, andrHeight-5, andrHeight+25, transferTx.Packet)
	require.NoError(t, err)

	// Second time, this will make the counter == 1 since the account is now created.
	transferTx, err = chainA.SendIBCTransfer(ctx, channel.ChannelID, user1.KeyName(), transfer, memo)
	require.NoError(t, err)
	andrHeight, err = chainA.Height(ctx)
	require.NoError(t, err)

	_, err = testutil.PollForAck(ctx, chainA, andrHeight-5, andrHeight+25, transferTx.Packet)
	require.NoError(t, err)

	// Get the address on the other chain's side
	addr := helpers.GetIBCHooksUserAddress(t, ctx, chainA, channel.ChannelID, user1Addr)
	require.NotEmpty(t, addr)

	// Get funds on the receiving chain
	funds := helpers.GetIBCHookTotalFunds(t, ctx, chainB, contractAddr, addr)
	require.Equal(t, int(1), len(funds.Data.TotalFunds))

	var ibcDenom string
	for _, coin := range funds.Data.TotalFunds {
		if strings.HasPrefix(coin.Denom, "ibc/") {
			ibcDenom = coin.Denom
			break
		}
	}
	require.NotEmpty(t, ibcDenom)

	// ensure the count also increased to 1 as expected.
	count := helpers.GetIBCHookCount(t, ctx, chainB, contractAddr, addr)
	require.Equal(t, int64(1), count.Data.Count)

}
