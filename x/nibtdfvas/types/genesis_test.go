package types_test

import (
	"testing"

	// "andromedad/x/nibtdfvas/types"

	"github.com/andromedaprotocol/andromedad/x/nibtdfvas/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "invalid genesis state",
			genState: &types.GenesisState{
				Params: types.Params{
					TokenOutflowPerBlock: 5,
					DirectToValidatorPercent: 12,
				},
			},
			valid: false,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{
				Params: types.Params{
					TokenOutflowPerBlock: 3,
					DirectToValidatorPercent: 20,
				},
			},
			valid: true,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
