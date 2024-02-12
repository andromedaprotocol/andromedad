// package types

// import (
// 	"fmt"
// 	"strconv"

// 	"gopkg.in/yaml.v2"
// )

// var (
// 	KeyTxNibtdfvasPercent = []byte("TxNibtdfvasPercent")
// 	// TODO: Determine the default value
// 	DefaultTxNibtdfvasPercent = "0"
// )

// // NewParams creates a new Params instance
// func NewParams(
// 	txNibtdfvasPercent string,
// ) Params {
// 	return Params{
// 		TxNibtdfvasPercent: txNibtdfvasPercent,
// 	}
// }

// // DefaultParams returns a default set of parameters
// func DefaultParams() Params {
// 	return NewParams(
// 		DefaultTxNibtdfvasPercent,
// 	)
// }

// // Validate validates the set of params
// func (p Params) Validate() error {
// 	return validateTxNibtdfvasPercent(p.TxNibtdfvasPercent)
// }

// // String implements the Stringer interface.
// func (p Params) String() string {
// 	out, _ := yaml.Marshal(p)
// 	return string(out)
// }

// // validateTxNibtdfvasPercent validates the TxNibtdfvasPercent param
// func validateTxNibtdfvasPercent(v interface{}) error {
// 	txNibtdfvasPercent, ok := v.(string)
// 	if !ok {
// 		return fmt.Errorf("invalid parameter type: %T", v)
// 	}

// 	txNibtdfvasPercentInt, err := strconv.Atoi(txNibtdfvasPercent)
// 	if err != nil {
// 		return err
// 	}
// 	if txNibtdfvasPercentInt < 0 || txNibtdfvasPercentInt > 100 {
// 		return fmt.Errorf("fee must be between 0 and 100")
// 	}

// 	return nil
// }


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

	_, ok1 := p1.(int64)
	_, ok2 := p2.(int64)
	
	if !ok1 {
		return fmt.Errorf("TokenOutflowPerBlock must be int64")
	}

	if !ok2 {
		return fmt.Errorf("DirectToValidatorPercent must be int64")
	}

	// if TokenOutflowPerBlock != 30 && TokenOutflowPerBlock != 0 {
	// 	return fmt.Errorf("TokenOutflowPerBlock must be 3")
	// }

	// if DirectToValidatorPercent != 20 && DirectToValidatorPercent != 0 {
	// 	return fmt.Errorf("DirectToValidatorPercent must be 20")
	// }

	return nil
}



