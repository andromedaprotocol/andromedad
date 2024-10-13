package types

import (
	"fmt"
	"strconv"
)

var (
	KeyTxFeeBurnPercent = []byte("TxFeeBurnPercent")
	// TODO: Determine the default value
	DefaultTxFeeBurnPercent = "0"
)

// NewParams creates a new Params instance
func NewParams(
	txFeeBurnPercent string,
) Params {
	return Params{
		TxFeeBurnPercent: txFeeBurnPercent,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultTxFeeBurnPercent,
	)
}

// Validate validates the set of params
func (p Params) Validate() error {
	return validateTxFeeBurnPercent(p.TxFeeBurnPercent)
}

// validateTxFeeBurnPercent validates the TxFeeBurnPercent param
func validateTxFeeBurnPercent(v interface{}) error {
	txFeeBurnPercent, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	txFeeBurnPercentInt, err := strconv.Atoi(txFeeBurnPercent)
	if err != nil {
		return err
	}
	if txFeeBurnPercentInt < 0 || txFeeBurnPercentInt > 100 {
		return fmt.Errorf("fee must be between 0 and 100")
	}

	return nil
}
