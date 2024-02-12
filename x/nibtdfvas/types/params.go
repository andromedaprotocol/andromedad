package types

import (
	"fmt"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}

// Validate validates the set of params
func (p *Params) Validate() error {
	return validateGenesisParams(p.TokenOutflowPerBlock, p.DirectToValidatorPercent)
}


func validateGenesisParams(p1 interface{}, p2 interface{}) error {

	TokenOutflowPerBlock, ok1 := p1.(int64)
	DirectToValidatorPercent, ok2 := p2.(int64)
	
	if !ok1 {
		return fmt.Errorf("TokenOutflowPerBlock must be int64")
	}

	if !ok2 {
		return fmt.Errorf("DirectToValidatorPercent must be int64")
	}

	if TokenOutflowPerBlock != 30 && TokenOutflowPerBlock != 0 {
		return fmt.Errorf("TokenOutflowPerBlock must be 3")
	}

	if DirectToValidatorPercent != 20 && DirectToValidatorPercent != 0 {
		return fmt.Errorf("DirectToValidatorPercent must be 20")
	}

	return nil
}



