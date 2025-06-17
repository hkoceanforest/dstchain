package keeper

import (
	sdkmath "cosmossdk.io/math"
	"errors"
	"freemasonry.cc/blockchain/contracts"
	chainCore "freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/x/contract"
	"freemasonry.cc/blockchain/x/dao/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"math/big"
	"strings"
)

var _ evmtypes.EvmHooks = EVMHooks{}

type EVMHooks struct {
	k Keeper
}

func (k Keeper) EVMHooks() EVMHooks {
	return EVMHooks{k}
}

func (h EVMHooks) PostTxProcessing(ctx sdk.Context, msg core.Message, receipt *ethtypes.Receipt) error {
	return h.k.PostTxProcessing(ctx, msg, receipt)
}

func (k Keeper) PostTxProcessing(
	ctx sdk.Context,
	msg core.Message,
	receipt *ethtypes.Receipt,
) error {
	
	err := genesisIdoContract(ctx, k, msg, receipt)
	if err != nil {
		return err
	}
	err = authContract(ctx, k, msg, receipt)
	if err != nil {
		return err
	}
	err = redPacket(ctx, k, msg, receipt)
	if err != nil {
		return err
	}
	err = erc20Reward(ctx, k, msg, receipt)
	if err != nil {
		return err
	}
	err = swap(ctx, k, msg)
	if err != nil {
		return err
	}
	return nil
}

func swap(ctx sdk.Context, k Keeper, msg core.Message) error {
	logs := chainCore.BuildLog(chainCore.GetFuncName(), chainCore.LmChainDaoEvmHook)
	if msg.To() == nil {
		return nil
	}
	
	if strings.ToLower(msg.To().String()) != strings.ToLower(contract.ExchangeRouterContract.String()) {
		return nil
	}
	if !k.GetGenesisIdoEndMark(ctx) {
		logs.Warning("PostTxProcessing swap contract not yet open")
		return errors.New("swap contract not yet open")
	}
	return nil
}

func erc20Reward(ctx sdk.Context, k Keeper, msg core.Message, receipt *ethtypes.Receipt) error {
	logs := chainCore.BuildLog(chainCore.GetFuncName(), chainCore.LmChainDaoEvmHook)
	if msg.To() == nil {
		return nil
	}
	for _, log := range receipt.Logs {
		
		if len(log.Topics) != 3 {
			continue
		}
		erc20RewardAccount := common.BytesToAddress(chainCore.Erc20Reward.Bytes())
		
		toAddress := common.BytesToAddress(log.Topics[2].Bytes())
		
		if toAddress.String() != erc20RewardAccount.String() {
			continue
		}
		erc20Abi := contracts.UsdtContract.ABI
		lpAbi := contracts.LPContract.ABI
		eventID := log.Topics[0]
		event, err := erc20Abi.EventByID(eventID)
		if err != nil {
			continue
		}
		if event.Name == "Transfer" {
			erc20Event, err := erc20Abi.Unpack(event.Name, log.Data)
			if err != nil {
				logs.WithError(err).Error("failed to unpack erc20Event event ")
				return err
			}
			if len(erc20Event) == 0 {
				continue
			}
			amount := erc20Event[0].(*big.Int)
			
			resp, err := k.contractKeeper.CallEVM(ctx, lpAbi, erc20RewardAccount, log.Address, true, "token0")
			if err != nil {
				return err
			}
			unPackData0, err := lpAbi.Unpack("token0", resp.Ret)
			if err != nil {
				return err
			}
			token0Address := unPackData0[0].(common.Address)

			resp, err = k.contractKeeper.CallEVM(ctx, erc20Abi, erc20RewardAccount, token0Address, true, "symbol")
			if err != nil {
				return err
			}

			unPackSymbolData0, err := erc20Abi.Unpack("symbol", resp.Ret)
			if err != nil {
				return err
			}
			symbol0 := unPackSymbolData0[0].(string)

			resp, err = k.contractKeeper.CallEVM(ctx, lpAbi, erc20RewardAccount, log.Address, true, "token1")
			if err != nil {
				return err
			}
			unPackData1, err := lpAbi.Unpack("token1", resp.Ret)
			if err != nil {
				return err
			}
			token1Address := unPackData1[0].(common.Address)

			resp, err = k.contractKeeper.CallEVM(ctx, erc20Abi, erc20RewardAccount, token1Address, true, "symbol")
			if err != nil {
				return err
			}
			unPackSymbolData1, err := erc20Abi.Unpack("symbol", resp.Ret)
			if err != nil {
				return err
			}
			symbol1 := unPackSymbolData1[0].(string)
			
			erc20Swap := types.Erc20Swap{
				Contract: log.Address.String(),
				Symbol:   symbol0 + "-" + symbol1,
				Token0:   types.TokenInfo{Symbol: symbol0, Contract: token0Address.String()},
				Token1:   types.TokenInfo{Symbol: symbol1, Contract: token1Address.String()},
			}
			err = k.SetErc20Swap(ctx, log.Address.String(), erc20Swap)
			if err != nil {
				return err
			}
			
			err = k.AddErc20Reward(ctx, sdk.NewDecFromBigInt(amount), log.Address.String())
			if err != nil {
				return err
			}
		}
	}
	return nil
}


func authContract(ctx sdk.Context, k Keeper, msg core.Message, receipt *ethtypes.Receipt) error {
	logs := chainCore.BuildLog(chainCore.GetFuncName(), chainCore.LmChainDaoEvmHook)
	
	if msg.To() == nil {
		return nil
	}
	auth := contracts.AuthJSONContract.ABI
	for _, log := range receipt.Logs {
		
		if log.Address.String() != contract.AuthContract.String() {
			continue
		}
		if len(log.Topics) == 0 {
			continue
		}
		eventID := log.Topics[0]
		event, err := auth.EventByID(eventID)
		if err != nil {
			return err
		}
		if event.Name == types.AuthEventApproveDvmLog {
			authEvent, err := auth.Unpack(event.Name, log.Data)
			if err != nil {
				logs.WithError(err).Error("failed to unpack auth event ")
				return err
			}
			if len(authEvent) == 0 {
				continue
			}
			
			address := common.BytesToAddress(log.Topics[1].Bytes())
			addr := sdk.AccAddress(address.Bytes())
			
			toAddress := common.BytesToAddress(log.Topics[2].Bytes())

			clusterId, ok := authEvent[0].(string)
			if !ok {
				return errors.New("clusterId is not string")
			}
			blockNumbers, ok := authEvent[1].(*big.Int)
			if !ok {
				return errors.New("blockNumbers is not *big.Int")
			}

			err = k.DvmApprove(ctx, toAddress.String(), clusterId, blockNumbers.String(), addr.String())
			if err != nil {
				logs.WithError(err).Error("failed to PersonDvmApprove ")
				return err
			}
		}

	}
	return nil
}


func genesisIdoContract(ctx sdk.Context, k Keeper, msg core.Message, receipt *ethtypes.Receipt) error {
	logs := chainCore.BuildLog(chainCore.GetFuncName(), chainCore.LmChainDaoEvmHook)
	
	if msg.To() == nil {
		return nil
	}

	genesisIdo := contracts.GenesisIdoNContract.ABI

	for _, log := range receipt.Logs {
		
		if log.Address.String() != contract.GenesisIdoContract.String() {
			continue
		}
		if len(log.Topics) == 0 {
			continue
		}
		eventID := log.Topics[0]
		event, err := genesisIdo.EventByID(eventID)
		if err != nil {
			return err
		}
		if event.Name == types.GenesisIdoEvent {

			genesisIdoEvent, err := genesisIdo.Unpack(event.Name, log.Data)
			if err != nil {
				logs.WithError(err).Error("failed to unpack transfer event ")
				return err
			}

			if len(genesisIdoEvent) == 0 {
				continue
			}

			address := common.BytesToAddress(log.Topics[1].Bytes())
			addr := sdk.AccAddress(address.Bytes())
			
			
			
			
			
			
			
			
			tokens, ok := genesisIdoEvent[1].(*big.Int)
			if !ok || tokens == nil || tokens.Sign() != 1 {
				continue
			}
			
			err = k.BurnGetNotActivePower(ctx, addr, sdkmath.NewIntFromBigInt(tokens))
			if err != nil {
				logs.WithError(err).Error("failed to Burn Get NotActivePower ")
				return err
			}
		}
	}
	return nil
}


func redPacket(ctx sdk.Context, k Keeper, msg core.Message, receipt *ethtypes.Receipt) error {
	logs := chainCore.BuildLog(chainCore.GetFuncName(), chainCore.LmChainDaoEvmHook)
	
	contractAddr := k.contractKeeper.GetRedPacketContractAddress(ctx)
	if msg.To() == nil || strings.ToLower(msg.To().String()) != strings.ToLower(contractAddr) {
		return nil
	}
	abi := contracts.RedPacketContract.ABI
	
	for _, log := range receipt.Logs {
		if len(log.Topics) == 0 {
			continue
		}
		eventID := log.Topics[0]
		event, err := abi.EventByID(eventID)
		if err != nil {
			continue
		}
		if event.Name == types.RedPacketLog {
			redPacketEvent, err := abi.Unpack(event.Name, log.Data)
			if err != nil {
				logs.WithError(err).Error("failed to unpack auth event ")
				return err
			}
			if len(redPacketEvent) == 0 {
				continue
			}
			
			rid := redPacketEvent[1].(*big.Int).String()
			
			clusterChatId := redPacketEvent[2].(string)

			clusterId, err := k.GetClusterId(ctx, clusterChatId)
			if err != nil {
				return err
			}

			
			err = k.SetRidClusterId(ctx, rid, clusterId)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
