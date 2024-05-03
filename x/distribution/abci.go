package distribution

import (
	"fmt"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"

	"github.com/andromedaprotocol/andromedad/x/distribution/keeper"
	"github.com/andromedaprotocol/andromedad/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

// BeginBlocker sets the proposer for determining distribution during endblock
// and distribute rewards for the previous block.
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	govKeeper := k.GovKeeper // Assuming you have access to the GovKeeper
	// Create a new proposal from v1.Proposal

	// Fetch proposal 1
	proposal, found := govKeeper.GetProposal(ctx, 1)
	if !found {
		panic(fmt.Sprintf("proposal %d does not exist", 1))
	}

	// Assuming 'Proposal' is a struct with fields that need to be copied.
	newProposal := v1.Proposal{
		Id:               6, // This should be a new unique ID
		Messages:         proposal.Messages,
		Status:           proposal.Status,
		FinalTallyResult: proposal.FinalTallyResult,
		SubmitTime:       proposal.SubmitTime,
		DepositEndTime:   proposal.DepositEndTime,
		TotalDeposit:     proposal.TotalDeposit,
		VotingStartTime:  proposal.VotingStartTime,
		VotingEndTime:    proposal.VotingEndTime,
		Metadata:         proposal.Metadata,
		Title:            proposal.Title,
		Summary:          proposal.Summary,
		Proposer:         proposal.Proposer,
	}

	govKeeper.SetProposal(ctx, newProposal)
	ctx.Logger().Info("Proposal changed", "proposalID", 6)
	govKeeper.RemoveFromInactiveProposalQueue(ctx, 6, *proposal.DepositEndTime)

	// determine the total power signing the block
	var previousTotalPower int64
	for _, voteInfo := range req.LastCommitInfo.GetVotes() {
		previousTotalPower += voteInfo.Validator.Power
	}

	// TODO this is Tendermint-dependent
	// ref https://github.com/cosmos/cosmos-sdk/issues/3095
	if ctx.BlockHeight() > 1 {
		k.AllocateTokens(ctx, previousTotalPower, req.LastCommitInfo.GetVotes())
	}

	// record the proposer for when we payout on the next block
	consAddr := sdk.ConsAddress(req.Header.ProposerAddress)
	k.SetPreviousProposerConsAddr(ctx, consAddr)
}
