package types

import (
	"fmt"
	"freemasonry.cc/blockchain/core"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	
	DefaultIndexNumHeight = int64(100)
	
	DefaultRedeemFeeHeight = int64(432000)
	
	DefaultRedeemFee = sdk.NewDec(1).Quo(sdk.NewDec(10))
	
	DefaultMinDelegate = sdk.NewCoin(core.GovDenom, sdk.NewInt(10).Mul(sdk.NewInt(core.RealToLedgerRateInt64)))
	
	DefaultValidity = int64(5256000)
)

var (
	KeyIndexNumHeight  = []byte("IndexNumHeight")
	KeyRedeemFeeHeight = []byte("RedeemFeeHeight")
	KeyRedeemFee       = []byte("RedeemFee")
	KeyMinDelegate     = []byte("MinDelegate")
	KeyValidity        = []byte("Validity")
)

func NewParams(
	IndexNumHeight int64,
	RedeemFeeHeight int64,
	RedeemFee sdk.Dec,
	MinDelegate sdk.Coin,
	Validity int64,
) Params {
	return Params{
		IndexNumHeight:  IndexNumHeight,
		RedeemFeeHeight: RedeemFeeHeight,
		RedeemFee:       RedeemFee,
		MinDelegate:     MinDelegate,
		Validity:        Validity,
	}
}

func DefaultParams() Params {

	return Params{
		IndexNumHeight:  DefaultIndexNumHeight,
		RedeemFeeHeight: DefaultRedeemFeeHeight,
		RedeemFee:       DefaultRedeemFee,
		MinDelegate:     DefaultMinDelegate,
		Validity:        DefaultValidity,
	}
}

func (p Params) Validate() error {
	if err := validateIndexNumHeight(p.IndexNumHeight); err != nil {
		return err
	}

	if err := validateRedeemFeeHeight(p.RedeemFeeHeight); err != nil {
		return err
	}

	if err := validateRedeemFee(p.RedeemFee); err != nil {
		return err
	}

	if err := validateMinDelegate(p.MinDelegate); err != nil {
		return err
	}
	if err := validateValidity(p.Validity); err != nil {
		return err
	}
	return nil
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyIndexNumHeight, &p.IndexNumHeight, validateIndexNumHeight),
		paramtypes.NewParamSetPair(KeyRedeemFeeHeight, &p.RedeemFeeHeight, validateRedeemFeeHeight),
		paramtypes.NewParamSetPair(KeyRedeemFee, &p.RedeemFee, validateRedeemFee),
		paramtypes.NewParamSetPair(KeyMinDelegate, &p.MinDelegate, validateMinDelegate),
		paramtypes.NewParamSetPair(KeyValidity, &p.Validity, validateValidity),
	}
}

func validateIndexNumHeight(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v <= 0 {
		return fmt.Errorf("IndexNumHeight must be positive: %d", v)
	}
	return nil
}

func validateRedeemFeeHeight(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v <= 0 {
		return fmt.Errorf("RedeemFeeHeight must be positive: %d", v)
	}
	return nil
}

func validateRedeemFee(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("RedeemFee cannot be negative: %s", v)
	}

	return nil
}

func validateMinDelegate(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.Amount.IsNegative() {
		return fmt.Errorf("MinDelegate cannot be negative: %s", v)
	}
	return nil
}

func validateValidity(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v <= 0 {
		return fmt.Errorf("Validity must be positive: %d", v)
	}
	return nil
}

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable(
		paramtypes.NewParamSetPair(KeyIndexNumHeight, DefaultParams().IndexNumHeight, validateIndexNumHeight),
		paramtypes.NewParamSetPair(KeyRedeemFeeHeight, DefaultParams().RedeemFeeHeight, validateRedeemFeeHeight),
		paramtypes.NewParamSetPair(KeyRedeemFee, DefaultParams().RedeemFee, validateRedeemFee),
		paramtypes.NewParamSetPair(KeyMinDelegate, DefaultParams().MinDelegate, validateMinDelegate),
		paramtypes.NewParamSetPair(KeyValidity, DefaultParams().Validity, validateValidity),
	)
}
