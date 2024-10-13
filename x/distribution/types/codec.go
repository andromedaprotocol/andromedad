package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

// RegisterLegacyAminoCodec registers the necessary x/distribution interfaces
// and concrete types on the provided LegacyAmino codec. These types are used
// for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgWithdrawDelegatorReward{}, "andromedad/MsgWithdrawDelegationReward")
	legacy.RegisterAminoMsg(cdc, &MsgWithdrawValidatorCommission{}, "andromedad/MsgWithdrawValCommission")
	legacy.RegisterAminoMsg(cdc, &MsgSetWithdrawAddress{}, "andromedad/MsgModifyWithdrawAddress")
	legacy.RegisterAminoMsg(cdc, &MsgFundCommunityPool{}, "andromedad/MsgFundCommunityPool")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "andromedad/distribution/MsgUpdateParams")
	legacy.RegisterAminoMsg(cdc, &MsgCommunityPoolSpend{}, "andromedad/distr/MsgCommunityPoolSpend")
	legacy.RegisterAminoMsg(cdc, &MsgDepositValidatorRewardsPool{}, "andromedad/distr/MsgDepositValRewards")
	legacy.RegisterAminoMsg(cdc, &MsgFundRewardsPool{}, "andromedad/MsgFundRewardsPool")

	cdc.RegisterConcrete(Params{}, "andromedad/x/distribution/Params", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgWithdrawDelegatorReward{},
		&MsgWithdrawValidatorCommission{},
		&MsgSetWithdrawAddress{},
		&MsgFundCommunityPool{},
		&MsgUpdateParams{},
		&MsgCommunityPoolSpend{},
		&MsgDepositValidatorRewardsPool{},
	)

	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&CommunityPoolSpendProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
