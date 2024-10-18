package v2

import (
	"context"
	"fmt"
	"os"
	"time"

	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	circuittypes "cosmossdk.io/x/circuit/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/andromedaprotocol/andromedad/app/upgrades"
)

const UpgradeName = "v2"

// NewUpgrade constructor
func NewUpgradeV2() upgrades.Upgrade {
	return upgrades.Upgrade{
		UpgradeName:          UpgradeName,
		CreateUpgradeHandler: CreateV2UpgradeHandler,
		StoreUpgrades: storetypes.StoreUpgrades{
			Added: []string{
				circuittypes.ModuleName,
			},
		},
	}
}

func CreateV2UpgradeHandler(
	mm upgrades.ModuleManager,
	configurator module.Configurator,
	ak *upgrades.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		logger := log.NewLogger(os.Stdout)
		currHeight := sdkCtx.BlockHeight()

		ak.ConsensusParamsKeeper.ParamsStore.Set(ctx, cmtproto.ConsensusParams{
			Block: &cmtproto.BlockParams{
				MaxBytes: 22020096,
				MaxGas:   -1,
			},
			Evidence: &cmtproto.EvidenceParams{
				MaxAgeNumBlocks: 100000,
				MaxAgeDuration:  time.Duration(48) * time.Hour,
				MaxBytes:        1048576,
			},
			Abci: &cmtproto.ABCIParams{
				VoteExtensionsEnableHeight: currHeight + 1,
			},
			Version: &cmtproto.VersionParams{},
			Validator: &cmtproto.ValidatorParams{
				PubKeyTypes: []string{"ed25519"},
			},
		})

		logger.Info(fmt.Sprintf("pre migrate version map: %v", fromVM))
		vm, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return nil, err
		}
		logger.Info(fmt.Sprintf("post migrate version map: %v", vm))

		// x/Staking - set minimum commission to 0.050000000000000000
		stakingParams, err := ak.StakingKeeper.GetParams(ctx)
		if err != nil {
			return nil, err
		}
		stakingParams.MinCommissionRate = sdkmath.LegacyNewDecWithPrec(5, 2)
		if err := ak.StakingKeeper.SetParams(ctx, stakingParams); err != nil {
			return nil, err
		}

		return vm, nil
	}
}
