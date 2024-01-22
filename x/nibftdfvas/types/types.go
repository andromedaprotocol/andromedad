package types

import (
	"github.com/cosmos/cosmos-sdk/types"
)

type NibtdfvasParams struct {
	TokenOutflowPerBlock int			// `json: "token_outfow_per_block"`
	DirectToValidatorPercent int	// `json:"direct_to_validator_percent"`
}

type NibtdfvasAccount struct {
	Address string
	Balance int
}

type MsgTokenDistribution struct {
	Sender string
	Amount int
}

func NewMsgTokenDistribution(sender string, amount int) MsgTokenDistribution {
	return MsgTokenDistribution{
		Sender: sender,
		Amount: amount,
	}
}

func (msg MsgTokenDistribution) Route() string {
	return "nibtdfvas"
}

func (msg MsgTokenDistribution) Type() string {
	return "token_distribution"
}

func (msg MsgTokenDistribution) ValidateBasic() error {
	return nil
}

func (msg MsgTokenDistribution) GetSignBytes() []byte {
	return types.MustSortJson(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgTokenDistribution) GetSigners() []types.AccAddress {
	sender, err := types.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []types.AccAddress{sender}
}