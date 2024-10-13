package testutil

import (
	_ "github.com/andromedaprotocol/andromedad/x/distribution" // import as blank for app wiring
	_ "github.com/cosmos/cosmos-sdk/x/auth"                    // import as blank for app wiring
	_ "github.com/cosmos/cosmos-sdk/x/auth/tx/config"          // import as blank for app wiring
	_ "github.com/cosmos/cosmos-sdk/x/bank"                    // import as blank for app wiring
	_ "github.com/cosmos/cosmos-sdk/x/consensus"               // import as blank for app wiring
	_ "github.com/cosmos/cosmos-sdk/x/genutil"                 // import as blank for app wiring
	_ "github.com/cosmos/cosmos-sdk/x/mint"                    // import as blank for app wiring
	_ "github.com/cosmos/cosmos-sdk/x/params"                  // import as blank for app wiring
	_ "github.com/cosmos/cosmos-sdk/x/staking"                 // import as blank for app wiring

	appv1alpha1 "cosmossdk.io/api/cosmos/app/v1alpha1"
	authmodulev1 "cosmossdk.io/api/cosmos/auth/module/v1"
	"cosmossdk.io/core/appconfig"

	"github.com/cosmos/cosmos-sdk/testutil/configurator"

	distrtypes "github.com/andromedaprotocol/andromedad/x/distribution/types"
)

var AppConfig = configurator.NewAppConfig(
	// configurator.AuthModule(),
	buildAuth(),
	configurator.BankModule(),
	configurator.StakingModule(),
	configurator.TxModule(),
	configurator.ConsensusModule(),
	configurator.ParamsModule(),
	configurator.GenutilModule(),
	configurator.DistributionModule(),
	configurator.MintModule(),
)

func buildAuth() configurator.ModuleOption {
	return func(config *configurator.Config) {
		config.ModuleConfigs["auth"] = &appv1alpha1.ModuleConfig{
			Name: "auth",
			Config: appconfig.WrapAny(&authmodulev1.Module{
				Bech32Prefix: "cosmos",
				ModuleAccountPermissions: []*authmodulev1.ModuleAccountPermission{
					{Account: "fee_collector", Permissions: []string{"burner"}},
					{Account: "distribution"},
					{Account: "mint", Permissions: []string{"minter"}},
					{Account: "bonded_tokens_pool", Permissions: []string{"burner", "staking"}},
					{Account: "not_bonded_tokens_pool", Permissions: []string{"burner", "staking"}},
					{Account: "gov", Permissions: []string{"burner"}},
					{Account: "nft"},
					// Andromeda app
					{Account: distrtypes.RewardsDripperName},
				},
			}),
		}
	}
}
