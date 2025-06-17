package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	
	
	MainnetMinGasPrices = sdk.NewDec(20_000_000_000)
	
	
	MainnetMinGasMultiplier = sdk.NewDecWithPrec(5, 1)
)
