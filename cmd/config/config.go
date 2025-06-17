package config

import (
	"freemasonry.cc/blockchain/core"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/evmos/ethermint/types"
)

const (
	
	Bech32Prefix = "dst"

	
	Bech32PrefixAccAddr = Bech32Prefix
	
	Bech32PrefixAccPub = Bech32Prefix + sdk.PrefixPublic
	
	Bech32PrefixValAddr = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator
	
	Bech32PrefixValPub = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator + sdk.PrefixPublic
	
	Bech32PrefixConsAddr = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus
	
	Bech32PrefixConsPub = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus + sdk.PrefixPublic
)

func SetBech32Prefixes(config *sdk.Config) {
	config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)
}

func SetBip44CoinType(config *sdk.Config) {
	config.SetCoinType(types.Bip44CoinType)
	config.SetPurpose(sdk.Purpose)                  
	config.SetFullFundraiserPath(types.BIP44HDPath) 
}

func RegisterDenoms() {
	
	
	

	if err := sdk.RegisterDenom(core.BaseDenom, sdk.NewDecWithPrec(1, types.BaseDenomUnit)); err != nil {
		panic(err)
	}
}
