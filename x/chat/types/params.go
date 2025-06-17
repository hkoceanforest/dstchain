package types

import (
	"freemasonry.cc/blockchain/core"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	KeyMaxPhoneNumber         = []byte("MaxPhoneNumber")
	KeyDestroyPhoneNumberCoin = []byte("DestroyPhoneNumberCoin")
)

func NewParams(
	maxPhoneNumber uint64,
	destroyPhoneNumberCoin sdk.Coin,
) Params {
	return Params{
		MaxPhoneNumber:         maxPhoneNumber,
		DestroyPhoneNumberCoin: destroyPhoneNumberCoin,
	}
}

func DefaultParams() Params {
	defalultDestroyPhoneNumberCoinInt, ok := sdk.NewIntFromString("10000000000000000000")
	if !ok {
		panic("Params DestroyPhoneNumberCoin NewIntFromString Error")
	}

	return Params{
		MaxPhoneNumber:         10,
		DestroyPhoneNumberCoin: sdk.NewCoin(core.BaseDenom, defalultDestroyPhoneNumberCoinInt),
	}
}

func (p Params) Validate() error {

	if err := validateMaxPhoneNumber(p.MaxPhoneNumber); err != nil {
		return err
	}

	if err := validateDestroy(p.DestroyPhoneNumberCoin); err != nil {
		return err
	}

	return nil
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMaxPhoneNumber, &p.MaxPhoneNumber, validateMaxPhoneNumber),
		paramtypes.NewParamSetPair(KeyDestroyPhoneNumberCoin, &p.DestroyPhoneNumberCoin, validateDestroy),
	}
}

func validateMaxPhoneNumber(i interface{}) error {
	
	maxPhoneNumber := i.(uint64)
	if maxPhoneNumber < 1 || maxPhoneNumber > 10000 {
		return core.ErrMaxPhoneNumber
	}

	return nil
}

func validateDestroy(i interface{}) error {
	
	destroyCoin := i.(sdk.Coin)
	if destroyCoin.Denom != core.BaseDenom {
		return core.ErrDestroyCoinDenom
	}

	if destroyCoin.Amount.LT(sdk.NewInt(1000000000000000000)) {
		return core.ErrDestroyCoinDenom
	}
	return nil
}

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable(
		paramtypes.NewParamSetPair(KeyMaxPhoneNumber, DefaultParams().MaxPhoneNumber, validateMaxPhoneNumber),
		paramtypes.NewParamSetPair(KeyDestroyPhoneNumberCoin, DefaultParams().DestroyPhoneNumberCoin, validateDestroy),
	)
}
