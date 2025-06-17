package types

import (
	"fmt"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	
	DefaultDays = int64(8640000)
)

var (
	KeyDays = []byte("Days")
)

func NewParams(
	days int64,
) Params {
	return Params{}
}

func DefaultParams() Params {

	return Params{}
}

func (p Params) Validate() error {
	
	
	
	return nil
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	
	
	
	return nil
}

func validateDays(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v <= 0 {
		return fmt.Errorf("Days must be positive: %d", v)
	}
	return nil
}

func ParamKeyTable() paramtypes.KeyTable {
	
	
	
	return paramtypes.KeyTable{}
}
