package ante

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	errorsmod "cosmossdk.io/errors"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	vestingtypes "github.com/evmos/evmos/v10/x/vesting/types"
)

type EthVestingTransactionDecorator struct {
	ak evmtypes.AccountKeeper
}

func NewEthVestingTransactionDecorator(ak evmtypes.AccountKeeper) EthVestingTransactionDecorator {
	return EthVestingTransactionDecorator{
		ak: ak,
	}
}


func (vtd EthVestingTransactionDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	for _, msg := range tx.GetMsgs() {
		msgEthTx, ok := msg.(*evmtypes.MsgEthereumTx)
		if !ok {
			return ctx, errorsmod.Wrapf(errortypes.ErrUnknownRequest,
				"invalid message type %T, expected %T", msg, (*evmtypes.MsgEthereumTx)(nil),
			)
		}

		acc := vtd.ak.GetAccount(ctx, msgEthTx.GetFrom())
		if acc == nil {
			return ctx, errorsmod.Wrapf(errortypes.ErrUnknownAddress,
				"account %s does not exist", acc)
		}

		
		clawbackAccount, isClawback := acc.(*vestingtypes.ClawbackVestingAccount)
		if !isClawback {
			return next(ctx, tx, simulate)
		}

		
		
		
		vested := clawbackAccount.GetVestedCoins(ctx.BlockTime())
		if len(vested) == 0 {
			return ctx, errorsmod.Wrapf(vestingtypes.ErrInsufficientVestedCoins,
				"cannot perform Ethereum tx with clawback vesting account, that has no vested coins: %s", vested,
			)
		}

		
		islocked := clawbackAccount.HasLockedCoins(ctx.BlockTime())
		if islocked {
			return ctx, errorsmod.Wrapf(vestingtypes.ErrVestingLockup,
				"cannot perform Ethereum tx with clawback vesting account, that has locked coins: %s", vested,
			)
		}
	}

	return next(ctx, tx, simulate)
}


type VestingDelegationDecorator struct {
	ak  evmtypes.AccountKeeper
	sk  vestingtypes.StakingKeeper
	cdc codec.BinaryCodec
}

func NewVestingDelegationDecorator(ak evmtypes.AccountKeeper, sk vestingtypes.StakingKeeper, cdc codec.BinaryCodec) VestingDelegationDecorator {
	return VestingDelegationDecorator{
		ak:  ak,
		sk:  sk,
		cdc: cdc,
	}
}

func (vdd VestingDelegationDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	for _, msg := range tx.GetMsgs() {
		switch msg := msg.(type) {
		case *authz.MsgExec:
			
			if err := vdd.validateAuthz(ctx, msg); err != nil {
				return ctx, err
			}
		default:
			if err := vdd.validateMsg(ctx, msg); err != nil {
				return ctx, err
			}
		}
	}

	return next(ctx, tx, simulate)
}

func (vdd VestingDelegationDecorator) validateAuthz(ctx sdk.Context, execMsg *authz.MsgExec) error {
	for _, v := range execMsg.Msgs {
		var innerMsg sdk.Msg
		if err := vdd.cdc.UnpackAny(v, &innerMsg); err != nil {
			return errorsmod.Wrap(err, "cannot unmarshal authz exec msgs")
		}

		if err := vdd.validateMsg(ctx, innerMsg); err != nil {
			return err
		}
	}

	return nil
}

func (vdd VestingDelegationDecorator) validateMsg(ctx sdk.Context, msg sdk.Msg) error {
	delegateMsg, ok := msg.(*stakingtypes.MsgDelegate)
	if !ok {
		return nil
	}

	for _, addr := range msg.GetSigners() {
		acc := vdd.ak.GetAccount(ctx, addr)
		if acc == nil {
			return errorsmod.Wrapf(
				errortypes.ErrUnknownAddress,
				"account %s does not exist", addr,
			)
		}

		clawbackAccount, isClawback := acc.(*vestingtypes.ClawbackVestingAccount)
		if !isClawback {
			
			return nil
		}

		
		bondDenom := vdd.sk.BondDenom(ctx)
		coins := clawbackAccount.GetVestedOnly(ctx.BlockTime())
		if coins == nil || coins.Empty() {
			return errorsmod.Wrap(
				vestingtypes.ErrInsufficientVestedCoins,
				"account has no vested coins",
			)
		}

		vested := coins.AmountOf(bondDenom)
		if vested.LT(delegateMsg.Amount.Amount) {
			return errorsmod.Wrapf(
				vestingtypes.ErrInsufficientVestedCoins,
				"cannot delegate unvested coins. coins vested < delegation amount (%s < %s)",
				vested, delegateMsg.Amount.Amount,
			)
		}
	}

	return nil
}
