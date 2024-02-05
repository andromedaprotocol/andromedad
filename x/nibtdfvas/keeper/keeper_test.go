package keeper_test
// import (
// 	"testing"
// 	// "fmt"

// 	// "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
// 	// "github.com/cosmos/cosmos-sdk/store/types"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/stretchr/testify/require"
// 	// "github.com/andromedaprotocol/andromedad/testutil/keeper"
// 	keepertest "github.com/andromedaprotocol/andromedad/testutil/keeper"
// 	"github.com/andromedaprotocol/andromedad/x/nibtdfvas/types"
// )

// func TestMint(t *testing.T) {
// 	// Create a mock context and a Keeper instance
// 	// ctx, keeper := createMockContextAndKeeper()
// 	keeper, ctx := keepertest.NibtdfvasKeeper(t)

// 	// Mint some coins
// 	mintCoins := sdk.NewCoins(sdk.NewCoin("token", sdk.NewInt(100)))


// 	// moduleAcc := keeper.GetNibtdfvasAccount(ctx)
// 	// if moduleAcc == nil {
// 	// 	panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
// 	// }


// 	err := keeper.Mint(ctx, types.ModuleName, mintCoins)


// 	// Verify that the minting was successful
// 	require.NoError(t, err)
// }

// func TestDistributeTokens(t *testing.T) {
// 	// Create a mock context and a Keeper instance
// 	ctx, keeper := createMockContextAndKeeper()

// 	// Distribute tokens
// 	keeper.DistributeTokens(ctx)
// }

// func createMockContextAndKeeper() (sdk.Context, *Keeper) {
// 	cdc := types.NewInterfaceRegistry()

// 	memKey := storetypes.NewMemoryStoreKey("memkey")
// 	storeKey := storetypes.NewMemoryStoreKey("storekey")

// 	ctx := sdk.NewContext(nil, sdk.Header{}, false, nil)
// 	authority := sdk.AccAddress(secp256k1.GenPrivKeqy().PubKey().Address())

// 	accountKeeper := // create a mock implementation of AccountKeeper
// 	bankKeeper :=    // create a mock implementation of BankKeeper
// 	stakingKeeper := // create a mock implementation of StakingKeeper

// 	keeper := NewKeeper(cdc, storeKey, memKey, authority, accountKeeper, bankKeeper, stakingKeeper)

// 	return ctx, keeper
// }