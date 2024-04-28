package v3

import (
	"github.com/andromedaprotocol/andromedad/x/distribution/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MigrateJSON accepts exported v2 (v0.46) x/distribution genesis state and migrates it to
// v3 (v0.47) x/distribution genesis state. The migration includes:
//
// Reset of the deprecated rewards to zero.
func MigrateJSON(oldState *types.GenesisState) *types.GenesisState {
	// reset deprecated rewards to zero
	oldState.Params.BaseProposerReward = sdk.ZeroDec()
	oldState.Params.BonusProposerReward = sdk.ZeroDec()

	return oldState
}
