package ante

import (
	"fmt"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/gogo/protobuf/proto"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authAnt "github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

type TxFeeChecker func(ctx sdk.Context, tx sdk.Tx) (sdk.Coins, int64, error)

type DeductFeeDecorator struct {
	accountKeeper  authAnt.AccountKeeper
	bankKeeper     evmtypes.BankKeeper
	feegrantKeeper authAnt.FeegrantKeeper
	txFeeChecker   TxFeeChecker
	daoKeeper      DaoKeeper
}

func NewDeductFeeDecorator(ak authAnt.AccountKeeper, bk evmtypes.BankKeeper, fk authAnt.FeegrantKeeper, tfc TxFeeChecker, dk DaoKeeper) DeductFeeDecorator {
	if tfc == nil {
		tfc = checkTxFeeWithValidatorMinGasPrices
	}

	return DeductFeeDecorator{
		accountKeeper:  ak,
		bankKeeper:     bk,
		feegrantKeeper: fk,
		txFeeChecker:   tfc,
		daoKeeper:      dk,
	}
}

func (dfd DeductFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	if !simulate && ctx.BlockHeight() > 0 && feeTx.GetGas() == 0 {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidGasLimit, "must provide positive gas")
	}

	var (
		priority int64
		err      error
	)

	fee := feeTx.GetFee()
	if !simulate {
		fee, priority, err = dfd.txFeeChecker(ctx, tx)
		if err != nil {
			return ctx, err
		}
	}
	if err := dfd.checkDeductFee(ctx, tx, fee); err != nil {
		return ctx, err
	}

	newCtx := ctx.WithPriority(priority)

	return next(newCtx, tx, simulate)
}

func (dfd DeductFeeDecorator) checkDeductFee(ctx sdk.Context, sdkTx sdk.Tx, fee sdk.Coins) error {
	feeTx, ok := sdkTx.(sdk.FeeTx)
	if !ok {
		return sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	if addr := dfd.accountKeeper.GetModuleAddress(types.FeeCollectorName); addr == nil {
		return fmt.Errorf("fee collector module account (%s) has not been set", types.FeeCollectorName)
	}

	feePayer := feeTx.FeePayer()

	feeGranter := feeTx.FeeGranter()
	deductFeesFrom := feePayer

	
	
	if feeGranter != nil {
		if dfd.feegrantKeeper == nil {
			return sdkerrors.ErrInvalidRequest.Wrap("fee grants are not enabled")
		} else if !feeGranter.Equals(feePayer) {
			err := dfd.feegrantKeeper.UseGrantedFees(ctx, feeGranter, feePayer, fee, sdkTx.GetMsgs())
			if err != nil {
				return sdkerrors.Wrapf(err, "%s does not not allow to pay fees for %s", feeGranter, feePayer)
			}
		}

		deductFeesFrom = feeGranter
	}

	deductFeesFromAcc := dfd.accountKeeper.GetAccount(ctx, deductFeesFrom)
	if deductFeesFromAcc == nil {
		deductFeesFromAcc = dfd.accountKeeper.NewAccountWithAddress(ctx, deductFeesFrom)
		dfd.accountKeeper.SetAccount(ctx, deductFeesFromAcc)
		
	}

	
	if !fee.IsZero() {
		var err error
		for _, msg := range sdkTx.GetMsgs() {
			if dfd.daoKeeper.FreeTxMsg(ctx, msg) { 
				payFee := fee
				if obj, oks := msg.(*bankTypes.MsgSend); oks {
					payFee, err = dfd.daoKeeper.CalculateSendFee(ctx, *obj, fee)
					if err != nil {
						return err
					}
				} else {
					
					payFee, err = dfd.daoKeeper.CalculateFee(ctx, msg, fee)
					if err != nil {
						return err
					}
				}
				if payFee != nil { 
					err = DeductFees(dfd.bankKeeper, ctx, deductFeesFromAcc, payFee)
					if err != nil {
						return err
					}
					events := sdk.Events{
						sdk.NewEvent(
							sdk.EventTypeTx,
							sdk.NewAttribute(sdk.AttributeKeyFee, payFee.String()),
							sdk.NewAttribute(sdk.AttributeKeyFeePayer, deductFeesFrom.String()),
							sdk.NewAttribute(sdk.AttributeKeySender, feePayer.String()),
						),
					}
					ctx.EventManager().EmitEvents(events)
					return nil
				}
			}
		}
	}
	return nil
}

func DeductFees(bankKeeper types.BankKeeper, ctx sdk.Context, acc types.AccountI, fees sdk.Coins) error {
	if !fees.IsValid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "invalid fee amount: %s", fees)
	}

	err := bankKeeper.SendCoinsFromAccountToModule(ctx, acc.GetAddress(), types.FeeCollectorName, fees)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, err.Error())
	}

	return nil
}

func ConcatMsgsType(msgs []sdk.Msg) string {
	msgsTypeString := ""

	for _, msg := range msgs {
		msgsTypeString += proto.MessageName(msg)

		msgsTypeString += ","
	}

	
	msgsTypeString = strings.TrimRight(msgsTypeString, ",")

	return msgsTypeString
}
