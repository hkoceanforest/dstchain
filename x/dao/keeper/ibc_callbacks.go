package keeper

import (
	sdkmath "cosmossdk.io/math"
	"errors"
	"freemasonry.cc/blockchain/contracts"
	"freemasonry.cc/blockchain/core"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v5/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v5/modules/core/04-channel/types"
	"github.com/cosmos/ibc-go/v5/modules/core/exported"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/evmos/v10/ibc"
	"math/big"
)

func (k Keeper) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	ack exported.Acknowledgement,
) exported.Acknowledgement {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	var data transfertypes.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return channeltypes.NewErrorAcknowledgement(err)
	}
	if !k.contractKeeper.IsERC20Enabled(ctx) {
		return ack
	}

	
	coin := ibc.GetReceivedCoin(
		packet.SourcePort, packet.SourceChannel,
		packet.DestinationPort, packet.DestinationChannel,
		data.Denom, data.Amount,
	)
	
	if coin.Denom == core.GovDenom {
		
		return ack
	}

	pairID := k.contractKeeper.GetTokenPairID(ctx, coin.Denom)
	if len(pairID) == 0 {
		logs.Warningf("received token %s is not registered", coin.Denom)
		
		
		return ack
	}

	pair, _ := k.contractKeeper.GetTokenPair(ctx, pairID)
	if !pair.Enabled {
		
		return ack
	}
	transferAmount, ok := sdk.NewIntFromString(data.Amount)
	if !ok {
		logs.WithError(core.ParseCoinError).Errorf("transfer amount %s", data.Amount)
		return channeltypes.NewErrorAcknowledgement(core.ParseCoinError)
	}
	if transferAmount.LTE(k.GetParams(ctx).CrossFee.FeeAmount) {
		logs.WithError(core.ErrIBCTransferAmount).Errorf("transfer amount is less than fee amount %s", transferAmount.String())
		return channeltypes.NewErrorAcknowledgement(core.ErrIBCTransferAmount)
	}
	
	err := k.DeductionToken(ctx, pair.Erc20Address, data.Receiver)
	if err != nil {
		logs.WithError(err).Errorf("failed to deduct token %s from %s", coin.Denom, data.Receiver)
		return channeltypes.NewErrorAcknowledgement(err)
	}
	return ack
}

func (k Keeper) DeductionToken(ctx sdk.Context, contractAddress, dstSender string) error {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	params := k.GetParams(ctx)
	feeAmount := params.CrossFee.FeeAmount
	feeAddress := common.HexToAddress(params.CrossFee.FeeCollectionAccount)
	dstAddr, err := sdk.AccAddressFromBech32(dstSender)
	if err != nil {
		return err
	}
	usdtABI := contracts.UsdtContract.ABI
	
	balanceResp, err := k.contractKeeper.CallEVM(ctx, usdtABI, common.BytesToAddress(dstAddr.Bytes()), common.HexToAddress(contractAddress), true, "balanceOf", common.BytesToAddress(dstAddr.Bytes()))
	if err != nil {
		logs.WithError(err).Error("failed to callEVM")
		return err
	}
	unPackData, err := usdtABI.Unpack("balanceOf", balanceResp.Ret)
	if err != nil {
		return err
	}
	balance := sdkmath.NewIntFromBigInt(unPackData[0].(*big.Int))
	if balance.LT(feeAmount) {
		return core.ErrInsufficientBalance
	}
	
	resp, err := k.contractKeeper.CallEVM(ctx, usdtABI, common.BytesToAddress(dstAddr.Bytes()), common.HexToAddress(contractAddress), true, "transfer", feeAddress, feeAmount.BigInt())
	if err != nil {
		logs.WithError(err).Error("failed to callEVM")
		return err
	}
	if resp.Failed() {
		logs.Warningf("failed to callEVM: %s", resp.VmError)
		return errors.New(resp.VmError)
	}
	return nil
}

func (k Keeper) OnAcknowledgementPacket(
	ctx sdk.Context, _ channeltypes.Packet,
	data transfertypes.FungibleTokenPacketData,
	ack channeltypes.Acknowledgement,
) error {
	switch ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:
		return k.RefundFeeFromPacket(ctx, data)
	default:
		
		
		return nil
	}
}
func (k Keeper) OnTimeoutPacket(ctx sdk.Context, _ channeltypes.Packet, data transfertypes.FungibleTokenPacketData) error {
	return k.RefundFeeFromPacket(ctx, data)
}

func (k Keeper) RefundFeeFromPacket(ctx sdk.Context, data transfertypes.FungibleTokenPacketData) error {
	
	return nil
}
