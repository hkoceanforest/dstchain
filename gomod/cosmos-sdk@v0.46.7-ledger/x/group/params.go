package group

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	DefaultDepositTokens = sdk.NewInt(1000000000000000000)
	KeyDeposit           = []byte("Deposit")
)

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(deposit sdk.Coins) Params {
	return Params{Deposit: deposit}
}

func DefaultParams() Params {
	return Params{
		Deposit: sdk.NewCoins(sdk.NewCoin("dst", DefaultDepositTokens)),
	}
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyDeposit, &p.Deposit, validateDeposit),
	}
}

func validateDeposit(i interface{}) error {
	v, ok := i.([]sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if !sdk.NewCoins(v...).IsValid() {
		return fmt.Errorf("invalid deposit: %s", v)
	}
	return nil
}
