package keeper

import (
	sdkmath "cosmossdk.io/math"
	"errors"
	"freemasonry.cc/blockchain/contracts"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/util"
	"freemasonry.cc/blockchain/x/dao/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

func (k Keeper) GetErc20Reward(ctx sdk.Context) (sdk.DecCoins, error) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetErc20RewardKey()
	token := sdk.DecCoins{}
	if store.Has(key) {
		bz := store.Get(key)
		err := util.Json.Unmarshal(bz, &token)
		if err != nil {
			return nil, err
		}
		return token, nil
	}
	return nil, nil
}

func (k Keeper) AddErc20Reward(ctx sdk.Context, amount sdk.Dec, daemon string) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetErc20RewardKey()
	token := sdk.DecCoins{}
	
	if daemon[:2] == "0x" {
		daemon = daemon[1:]
	}
	err := sdk.ValidateDenom(daemon)
	if err != nil {
		return err
	}
	newToken := sdk.NewDecCoinFromDec(daemon, amount)
	if store.Has(key) {
		bz := store.Get(key)
		err := util.Json.Unmarshal(bz, &token)
		if err != nil {
			return err
		}
	}
	token = token.Add(newToken)
	tokenByte, err := util.Json.Marshal(token)
	if err != nil {
		return err
	}
	store.Set(key, tokenByte)
	return nil
}

func (k Keeper) DeleteErc20Reward(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetErc20RewardKey()
	if store.Has(key) {
		store.Delete(key)
	}
}

func (k Keeper) SetErc20Swap(ctx sdk.Context, contract string, erc20Swap types.Erc20Swap) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetErc20SwapKey(contract)
	if !store.Has(key) {
		tokenByte, err := util.Json.Marshal(erc20Swap)
		if err != nil {
			return err
		}
		store.Set(key, tokenByte)
	}
	return nil
}

func (k Keeper) GetErc20Symbol(ctx sdk.Context, contract string) (*types.Erc20Swap, error) {
	store := ctx.KVStore(k.storeKey)
	
	if contract[:1] != "0" {
		contract = "0" + contract
	}
	key := types.GetErc20SwapKey(contract)
	token := new(types.Erc20Swap)
	if store.Has(key) {
		bz := store.Get(key)
		err := util.Json.Unmarshal(bz, token)
		if err != nil {
			return token, err
		}
		return token, nil
	}
	return nil, nil
}


func (k Keeper) SendErc20Reward(ctx sdk.Context, contract string, toAddress common.Address, amount sdkmath.Int) error {
	if contract[:1] != "0" {
		contract = "0" + contract
	}
	contractAddr := common.HexToAddress(contract)
	
	erc20Swap, err := k.GetErc20Symbol(ctx, contract)
	if err != nil {
		return err
	}
	
	if erc20Swap == nil || erc20Swap.Symbol == "" {
		return nil
	}
	erc20RewardAccount := common.BytesToAddress(core.Erc20Reward.Bytes())
	
	contractIsExist := k.contractKeeper.QueryContractIsExist(ctx, contractAddr)
	if !contractIsExist {
		return nil
	}
	lpAbi := contracts.LPContract.ABI
	
	resp, err := k.contractKeeper.CallEVM(ctx, lpAbi, erc20RewardAccount, contractAddr, true, "transfer", toAddress, amount.BigInt())
	if err != nil {
		return err
	}
	if resp.Failed() {
		return errors.New(resp.VmError)
	}
	return nil
}


func (k Keeper) GetErc20RewardCoins(historyReward sdk.Coins, finalReward sdk.Coins) sdk.Coins {
	var rewardCoins sdk.Coins
	if historyReward != nil && len(historyReward) > 0 {
		for _, coin := range historyReward {
			rewardCoins = rewardCoins.Add(coin)
		}
	}
	if finalReward != nil && len(finalReward) > 0 {
		for _, coin := range finalReward {
			if coin.Denom == core.BaseDenom {
				continue
			}
			rewardCoins = rewardCoins.Add(coin)
		}
	}
	return rewardCoins
}
