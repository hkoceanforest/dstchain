package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"freemasonry.cc/blockchain/core"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func ExportBank(jsonData json.RawMessage, totalStakeAmount sdk.Int) ([]byte, error) {
	bankState := bankTypes.GenesisState{}
	encCfg.Codec.MustUnmarshalJSON(jsonData, &bankState)

	
	
	var stakingPool, distPool, feePool, undbondedPool sdk.Int
	undbondedPool = sdk.ZeroInt()
	for _, balance := range bankState.Balances {
		if balance.Address == "dst1fl48vsnmsdzcv85q5d2q4z5ajdha8yu39kd84c" {
			stakingPool = totalStakeAmount
		}
		if balance.Address == "dst1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8hjq59s" {
			distPool = sdk.ZeroInt()
		}
		if balance.Address == "dst17xpfvakm2amg962yls6f84z3kell8c5lq58g8j" {
			feePool = balance.Coins.AmountOf(core.GovDenom)
		}
		if balance.Address == "dst1tygms3xhhs3yv487phx3dw4a95jn7t7l3k3krv" {
			undbondedPool = sdk.NewInt(2000000000000000000)
		}
	}

	nxnSupply := bankState.Supply.AmountOf(core.GovDenom)
	genesisAccountNxnAmount := nxnSupply.Sub(stakingPool).Sub(distPool).Sub(feePool).Sub(undbondedPool)

	fmt.Println("nxnSupply:", nxnSupply.String())
	fmt.Println("stakingPool:", stakingPool.String())
	fmt.Println("distPool:", distPool.String())
	fmt.Println("feePool:", feePool.String())
	fmt.Println("undbondedPool:", undbondedPool.String())
	fmt.Println("genesisAccountNxnAmount:", genesisAccountNxnAmount.String())

	bankState.Balances = make([]bankTypes.Balance, 0)

	genesisAccountDstAmount, ok := sdk.NewIntFromString("21600000000000000000000000")
	if !ok {
		return nil, errors.New("genesis account amount is invalid")
	}

	airdropAmount, ok := sdk.NewIntFromString("100000000000000000000000")
	if !ok {
		return nil, errors.New("airdrop amount is invalid")
	}

	daoAmount, ok := sdk.NewIntFromString("36000000000000000000000000")
	if !ok {
		return nil, errors.New("dao amount is invalid")
	}

	bankState.Balances = append(bankState.Balances, bankTypes.Balance{
		Address: "dst1ll30h0xykgduvxxfnpy4h6yzl0770pgneerr6w",
		Coins: sdk.NewCoins(
			sdk.NewCoin(core.BaseDenom, genesisAccountDstAmount),
			sdk.NewCoin(core.GovDenom, genesisAccountNxnAmount),
		),
	})
	
	bankState.Balances = append(bankState.Balances, bankTypes.Balance{
		Address: "dst1wtvnpd83e9lugyhk4n35z42jh5h8fqjmj5ejcl",
		Coins:   sdk.NewCoins(sdk.NewCoin(core.BaseDenom, airdropAmount)),
	})
	
	bankState.Balances = append(bankState.Balances, bankTypes.Balance{
		Address: "dst1vwr8z00ty7mqnk4dtchr9mn9j96nuh6wfn7scm",
		Coins:   sdk.NewCoins(sdk.NewCoin(core.BaseDenom, daoAmount)),
	})

	
	bankState.Balances = append(bankState.Balances, bankTypes.Balance{
		Address: "dst1fl48vsnmsdzcv85q5d2q4z5ajdha8yu39kd84c",
		Coins:   sdk.NewCoins(sdk.NewCoin(core.GovDenom, stakingPool)),
	})

	
	bankState.Balances = append(bankState.Balances, bankTypes.Balance{
		Address: "dst1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8hjq59s",
		Coins:   sdk.NewCoins(sdk.NewCoin(core.GovDenom, distPool)),
	})

	
	bankState.Balances = append(bankState.Balances, bankTypes.Balance{
		Address: "dst17xpfvakm2amg962yls6f84z3kell8c5lq58g8j",
		Coins:   sdk.NewCoins(sdk.NewCoin(core.GovDenom, feePool)),
	})

	
	bankState.Balances = append(bankState.Balances, bankTypes.Balance{
		Address: "dst1tygms3xhhs3yv487phx3dw4a95jn7t7l3k3krv",
		Coins:   sdk.NewCoins(sdk.NewCoin(core.GovDenom, undbondedPool)),
	})

	
	dstSupply, ok := sdk.NewIntFromString("57700000000000000000000000")
	bankState.Supply = sdk.NewCoins(
		sdk.NewCoin(core.GovDenom, nxnSupply),
		sdk.NewCoin(core.BaseDenom, dstSupply),
	)

	balancesByte := encCfg.Codec.MustMarshalJSON(&bankState)
	return balancesByte, nil
}
