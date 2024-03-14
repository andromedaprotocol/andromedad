package distribution

import (
	agithub.com/cosmos/cosmos-sdkosmos/cosmos-sdksmossdk.io/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	modulev1 "github.com/cosmos/cosmos-sdk/api/cosmos/distribution/module/v1"
	"github.com/cosmos/cosmos-sdk/core/appmodule"
	"github.com/cosmos/cosmos-sdk/depinject"
	"github.com/cosmos/cosmos-sdk/depinject/appconfig"

	"github.com/cosmos/cosmos-sdk/codec"
)

var _ depinject.OnePerModuleType = AppModule{}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() {}

func init() {
	appconfig.RegisterModule(&modulev1.Module{},
		appconfig.Provide(ProvideModule),
	)
}

type ModuleInputs struct {
	depinject.In

	Config      *modulev1.Module
	Environment appmodule.Environment
	Cdc         codec.Codec

	AccountKeeper types.AccountKeeper
	BankKeeper    types.BankKeeper
	StakingKeeper types.StakingKeeper
	PoolKeeper    types.PoolKeeper
}

type ModuleOutputs struct {
	depinject.Out

	DistrKeeper keeper.Keeper
	Module      appmodule.AppModule
	Hooks       staking.StakingHooksWrapper
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	feeCollectorName := in.Config.FeeCollectorName
	if feeCollectorName == "" {
		feeCollectorName = authtypes.FeeCollectorName
	}

	// default to governance authority if not provided
	authority := authtypes.NewModuleAddress(types.GovModuleName)
	if in.Config.Authority != "" {
		authority = authtypes.NewModuleAddressOrBech32Address(in.Config.Authority)
	}

	k := keeper.NewKeeper(
		in.Cdc,
		in.Environment,
		in.AccountKeeper,
		in.BankKeeper,
		in.StakingKeeper,
		in.PoolKeeper,
		feeCollectorName,
		authority.String(),
	)

	m := NewAppModule(in.Cdc, k, in.AccountKeeper, in.BankKeeper, in.StakingKeeper, in.PoolKeeper)

	return ModuleOutputs{
		DistrKeeper: k,
		Module:      m,
		Hooks:       staking.StakingHooksWrapper{StakingHooks: k.Hooks()},
	}
}
