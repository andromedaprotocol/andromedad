package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Verify interface at compile time
var (
	_ sdk.Msg = (*MsgSetWithdrawAddress)(nil)
	_ sdk.Msg = (*MsgWithdrawDelegatorReward)(nil)
	_ sdk.Msg = (*MsgWithdrawValidatorCommission)(nil)
	_ sdk.Msg = (*MsgUpdateParams)(nil)
	_ sdk.Msg = (*MsgCommunityPoolSpend)(nil)
	_ sdk.Msg = (*MsgDepositValidatorRewardsPool)(nil)
	_ sdk.Msg = (*MsgFundRewardsPool)(nil)
)

func NewMsgSetWithdrawAddress(delAddr, withdrawAddr sdk.AccAddress) *MsgSetWithdrawAddress {
	return &MsgSetWithdrawAddress{
		DelegatorAddress: delAddr.String(),
		WithdrawAddress:  withdrawAddr.String(),
	}
}

func NewMsgWithdrawDelegatorReward(delAddr, valAddr string) *MsgWithdrawDelegatorReward {
	return &MsgWithdrawDelegatorReward{
		DelegatorAddress: delAddr,
		ValidatorAddress: valAddr,
	}
}

func NewMsgWithdrawValidatorCommission(valAddr string) *MsgWithdrawValidatorCommission {
	return &MsgWithdrawValidatorCommission{
		ValidatorAddress: valAddr,
	}
}

// NewMsgFundCommunityPool returns a new MsgFundCommunityPool with a sender and
// a funding amount.
func NewMsgFundCommunityPool(amount sdk.Coins, depositor string) *MsgFundCommunityPool {
	return &MsgFundCommunityPool{
		Amount:    amount,
		Depositor: depositor,
	}
}

// NewMsgDepositValidatorRewardsPool returns a new MsgDepositValidatorRewardsPool
// with a depositor and a funding amount.
func NewMsgDepositValidatorRewardsPool(depositor, valAddr string, amount sdk.Coins) *MsgDepositValidatorRewardsPool {
	return &MsgDepositValidatorRewardsPool{
		Amount:           amount,
		Depositor:        depositor,
		ValidatorAddress: valAddr,
	}
}

// NewMsgFundRewardsPool returns a new MsgFundRewardsPool with a sender and
// a funding amount.
func NewMsgFundRewardsPool(amount sdk.Coins, depositor sdk.AccAddress) *MsgFundRewardsPool {
	return &MsgFundRewardsPool{
		Amount:    amount,
		Depositor: depositor.String(),
	}
}

// ValidateBasic performs basic MsgFundCommunityPool message validation.
func (msg MsgFundRewardsPool) ValidateBasic() error {
	if !msg.Amount.IsValid() {
		return fmt.Errorf("invalid amount: %s", msg.Amount)
	}
	if _, err := sdk.AccAddressFromBech32(msg.Depositor); err != nil {
		return fmt.Errorf("invalid depositor address: %s", err)
	}
	return nil
}
