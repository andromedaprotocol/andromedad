package types

import (
	"github.com/cosmos/cosmos-sdk/types"
)

type DAOParams struct {
	TokenOutflowPerBlock     int    `json:"token_outflow_per_block"`
	DirectToValidatorPercent int    `json:"direct_to_validator_percent"`
}

type GenesisState struct {
	DAOParams DAOParams `json:"dao_params"`
	// Add other genesis state as needed
}

type MyModuleAccount struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

type MsgTokenDistribution struct {
	Sender string `json:"sender"`
	Amount int    `json:"amount"`
}

var _ types.Msg = &MsgTokenDistribution{}

func NewMsgTokenDistribution(sender string, amount int) MsgTokenDistribution {
	return MsgTokenDistribution{
		Sender: sender,
		Amount: amount,
	}
}

func (msg MsgTokenDistribution) Route() string {
	return "mymodule"
}

func (msg MsgTokenDistribution) Type() string {
	return "token_distribution"
}

func (msg MsgTokenDistribution) ValidateBasic() error {
	// Add validation logic if needed
	return nil
}

func (msg MsgTokenDistribution) GetSignBytes() []byte {
	return types.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgTokenDistribution) GetSigners() []types.AccAddress {
	sender, err := types.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []types.AccAddress{sender}
}

func (msg MsgTokenDistribution) ValidateBasic() error {
	if msg.Amount <= 0 {
		return errors.New("amount should be greater than 0")
	}

	if _, err := types.AccAddressFromBech32(msg.Sender); err != nil{
		return errors.New("invalid sender address")
	}

	return nil
}


