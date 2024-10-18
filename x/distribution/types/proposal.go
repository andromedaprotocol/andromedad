package types

import (
	sdkdistrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

// GetTitle returns the title of a community pool spend proposal.
func (csp *CommunityPoolSpendProposal) GetTitle() string { return csp.Title }

// GetDescription returns the description of a community pool spend proposal.
func (csp *CommunityPoolSpendProposal) GetDescription() string { return csp.Description }

// GetDescription returns the routing key of a community pool spend proposal.
func (csp *CommunityPoolSpendProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a community pool spend proposal.
func (csp *CommunityPoolSpendProposal) ProposalType() string { return "CommunityPoolSpend" }

// ValidateBasic runs basic stateless validity checks
func (csp *CommunityPoolSpendProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(csp)
	if err != nil {
		return err
	}
	if !csp.Amount.IsValid() {
		return sdkdistrtypes.ErrInvalidProposalAmount
	}
	if csp.Recipient == "" {
		return sdkdistrtypes.ErrEmptyProposalRecipient
	}

	return nil
}
