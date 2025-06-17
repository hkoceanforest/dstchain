package keeper

import (
	"context"
	"encoding/base64"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/x/gateway/types"
	"github.com/armon/go-metrics"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/sirupsen/logrus"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"time"
)

var _ types.MsgServer = &msgServer{}

type msgServer struct {
	Keeper
	logPrefix string
}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper, logPrefix: "comm | msgServer | "}
}

func (k msgServer) GatewayUpload(goCtx context.Context, msg *types.MsgGatewayUpload) (*types.MsgEmptyResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)
	kvStore := ctx.KVStore(k.storeKey)

	accfrom, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil, core.ErrAddressFormat
	}

	valAddress := sdk.ValAddress(accfrom.Bytes())

	key := []byte(types.GatewayUploadKey + valAddress.String())

	kvStore.Set(key, msg.GatewayKeyInfo)

	return nil, nil
}

func (k msgServer) CreateSmartValidator(goCtx context.Context, msg *types.MsgCreateSmartValidator) (*types.MsgEmptyResponse, error) {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainMsgServer)
	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		log.WithError(err).Error("invalid delegator address")
		return &types.MsgEmptyResponse{}, err
	}
	valAddress := sdk.ValAddress(addr)
	_, found := k.stakingKeeper.GetValidator(ctx, valAddress)
	if found {
		return nil, stakingTypes.ErrValidatorOwnerExists
	}
	err = k.createValidator(ctx, addr, valAddress, *msg, msg.Value)
	if err != nil {
		log.WithError(err).Error("create validator failed")
		return &types.MsgEmptyResponse{}, err
	}
	return &types.MsgEmptyResponse{}, nil
}

func (k msgServer) GatewayRegister(goCtx context.Context, msg *types.MsgGatewayRegister) (*types.MsgEmptyResponse, error) {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainMsgServer)
	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		log.WithError(err).Error("invalid address")
		return &types.MsgEmptyResponse{}, err
	}
	if msg.Delegation == "" {
		msg.Delegation = "0"
	}
	delInt, ok := sdk.NewIntFromString(msg.Delegation)
	if !ok {
		log.Error("invalid delegation")
		return &types.MsgEmptyResponse{}, sdkerrors.Wrapf(
			core.ErrDelegationCoin, "invalid delegation : got %s", msg.Delegation)
	}
	params := k.GetParams(ctx)
	delegation := sdk.NewCoin(sdk.DefaultBondDenom, delInt)

	valAddress := sdk.ValAddress(addr)
	validator, found := k.stakingKeeper.GetValidator(ctx, valAddress)
	if !found {
		return nil, core.ErrValidatorNotFound
	}

	if !delegation.IsZero() {
		err = k.delegate(ctx, addr, valAddress, validator, delegation)
		if err != nil {
			log.WithError(err).Error("delegate failed")
			return nil, err
		}
	}

	delegate, found := k.stakingKeeper.GetDelegation(ctx, addr, valAddress)
	if !found {
		return nil, stakingTypes.ErrNoDelegation
	}

	if delegate.Shares.LT(sdk.NewDecFromInt(params.MinDelegate.Amount)) {
		return &types.MsgEmptyResponse{}, core.ErrGatewayDelegation
	}
	delegation = sdk.NewCoin(sdk.DefaultBondDenom, delegate.Shares.TruncateInt())

	err = k.SetGateway(ctx, *msg, delegation, valAddress.String())
	if err != nil {
		log.WithError(err).Error("set gateway failed")
		return &types.MsgEmptyResponse{}, err
	}
	machineAddr, err := sdk.AccAddressFromBech32(msg.MachineAddress)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		return &types.MsgEmptyResponse{}, err
	}
	accExists := k.accountKeeper.HasAccount(ctx, machineAddr)
	if !accExists {
		defer telemetry.IncrCounter(1, "new", "account")
		k.accountKeeper.SetAccount(ctx, k.accountKeeper.NewAccountWithAddress(ctx, machineAddr))
	}

	return &types.MsgEmptyResponse{}, nil
}

func (k msgServer) GatewayEdit(goCtx context.Context, msg *types.MsgGatewayEdit) (*types.MsgEmptyResponse, error) {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainMsgServer)
	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		log.WithError(err).Error("invalid address")
		return &types.MsgEmptyResponse{}, err
	}
	valAddress := sdk.ValAddress(addr)

	err = k.GatewayEdits(ctx, msg.GetGatewayUrl(), valAddress.String())
	if err != nil {
		log.WithError(err).Error("edit gateway failed")
		return &types.MsgEmptyResponse{}, err
	}
	return &types.MsgEmptyResponse{}, nil
}

func (k msgServer) GatewayIndexNum(goCtx context.Context, msg *types.MsgGatewayIndexNum) (*types.MsgEmptyResponse, error) {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainMsgServer)
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := msg.ValidateBasic()
	if err != nil {
		log.WithError(err).Error("validate basic failed")
		return nil, err
	}
	valAddr, valErr := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if valErr != nil {
		return nil, valErr
	}
	validator, found := k.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return nil, stakingTypes.ErrNoValidatorFound
	}
	delegatorAddress, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		log.WithError(err).Error("invalid delegator address")
		return nil, err
	}

	if validator.GetOperator().String() == sdk.ValAddress(delegatorAddress).String() {

		gateway, err := k.GetGatewayInfo(ctx, msg.ValidatorAddress)
		if err != nil {
			if err == core.ErrGatewayNotExist {
				return &types.MsgEmptyResponse{}, nil
			}
			log.WithError(err).Error("get gateway info failed")
			return nil, err
		}

		delegation, found := k.stakingKeeper.GetDelegation(ctx, delegatorAddress, valAddr)
		if !found {
			return nil, stakingTypes.ErrNoDelegation
		}
		params := k.GetParams(ctx)
		num := delegation.Shares.QuoInt(params.MinDelegate.Amount)
		gateway.GatewayQuota = num.TruncateInt64()
		if msg.IndexNumber != nil && len(msg.IndexNumber) > 0 {

			if num.Sub(sdk.NewDecFromInt(sdk.NewInt(int64(len(gateway.GatewayNum))))).LT(sdk.NewDec(int64(len(msg.IndexNumber)))) {
				log.WithFields(logrus.Fields{"num": num, "GatewayNum": len(gateway.GatewayNum), "IndexNumber": len(msg.IndexNumber)}).Error(core.ErrGatewayNum)
				return nil, core.ErrGatewayNum
			}
			indexNumArray, err := k.GatewayNumFilter(ctx, *gateway, msg.IndexNumber)
			if err != nil {
				log.WithError(err).Error("gateway num filter failed")
				return nil, err
			}

			err = k.SetGatewayNum(ctx, indexNumArray)
			if err != nil {
				log.WithError(err).Error("set gateway num failed")
				return nil, err
			}

			err = k.GatewayRedeemNumFilter(ctx, indexNumArray)
			if err != nil {
				return nil, err
			}
			gateway.GatewayNum = append(gateway.GatewayNum, indexNumArray...)

			gateway.Status = 0
		}

		err = k.UpdateGatewayInfo(ctx, *gateway)
		if err != nil {
			log.WithError(err).Error("update gateway info failed")
			return nil, err
		}
	}
	return &types.MsgEmptyResponse{}, nil
}

func (k msgServer) GatewayUndelegate(goCtx context.Context, msg *types.MsgGatewayUndelegate) (*types.MsgEmptyResponse, error) {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainMsgServer)
	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		log.WithError(err).Error("invalid validator address")
		return nil, err
	}
	delegatorAddress, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		log.WithError(err).Error("invalid delegator address")
		return nil, err
	}

	shares, err := k.stakingKeeper.ValidateUnbondAmount(
		ctx, delegatorAddress, addr, msg.Amount.Amount,
	)
	if err != nil {
		log.WithError(err).Error("validate unbond amount failed")
		return nil, err
	}
	bondDenom := k.stakingKeeper.BondDenom(ctx)
	if msg.Amount.Denom != bondDenom {
		log.WithFields(logrus.Fields{"bondDenom": bondDenom, "msg.Amount.Denom": msg.Amount.Denom}).Error("invalid bond denom")
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s, expected %s", msg.Amount.Denom, bondDenom,
		)
	}
	validator, found := k.stakingKeeper.GetValidator(ctx, addr)
	if !found {
		return nil, stakingTypes.ErrNoDelegatorForAddress
	}
	completionTime, returnAmount, err := k.Keeper.Undelegate(ctx, delegatorAddress, addr, validator, shares)
	if err != nil {
		log.WithError(err).Error("undelegate failed")
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			stakingTypes.EventTypeUnbond,
			sdk.NewAttribute(stakingTypes.AttributeKeyValidator, msg.ValidatorAddress),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyReturnAmount, returnAmount.String()),
			sdk.NewAttribute(stakingTypes.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),

			sdk.NewAttribute(stakingTypes.AttributeKeyDelegatorAddr, delegatorAddress.String()),
			sdk.NewAttribute(stakingTypes.AttributeKeyNewShares, shares.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, stakingTypes.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress),
		),
	})

	return &types.MsgEmptyResponse{}, nil
}

func (k msgServer) GatewayBeginRedelegate(goCtx context.Context, msg *types.MsgGatewayBeginRedelegate) (*types.MsgBeginRedelegateResponse, error) {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainMsgServer)
	ctx := sdk.UnwrapSDKContext(goCtx)
	valSrcAddr, err := sdk.ValAddressFromBech32(msg.ValidatorSrcAddress)
	if err != nil {
		log.WithError(err).Error("invalid source validator address")
		return nil, err
	}
	delegatorAddress, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		log.WithError(err).Error("invalid delegator address")
		return nil, err
	}

	shares, err := k.stakingKeeper.ValidateUnbondAmount(
		ctx, delegatorAddress, valSrcAddr, msg.Amount.Amount,
	)
	if err != nil {
		log.WithError(err).Error("validate unbond amount failed")
		return nil, err
	}

	bondDenom := k.stakingKeeper.BondDenom(ctx)
	if msg.Amount.Denom != bondDenom {
		log.WithFields(logrus.Fields{"bondDenom": bondDenom, "msg.Amount.Denom": msg.Amount.Denom}).Error("invalid bond denom")
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s, expected %s", msg.Amount.Denom, bondDenom,
		)
	}

	valDstAddr, err := sdk.ValAddressFromBech32(msg.ValidatorDstAddress)
	if err != nil {
		log.WithError(err).Error("invalid destination validator address")
		return nil, err
	}

	completionTime, err := k.stakingKeeper.BeginRedelegation(
		ctx, delegatorAddress, valSrcAddr, valDstAddr, shares,
	)
	if err != nil {
		log.WithError(err).Error("begin redelegation failed")
		return nil, err
	}

	if msg.Amount.Amount.IsInt64() {
		defer func() {
			telemetry.IncrCounter(1, types.ModuleName, "redelegate")
			telemetry.SetGaugeWithLabels(
				[]string{"tx", "msg", msg.Type()},
				float32(msg.Amount.Amount.Int64()),
				[]metrics.Label{telemetry.NewLabel("denom", msg.Amount.Denom)},
			)
		}()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			stakingTypes.EventTypeRedelegate,
			sdk.NewAttribute(stakingTypes.AttributeKeySrcValidator, msg.ValidatorSrcAddress),
			sdk.NewAttribute(stakingTypes.AttributeKeyDstValidator, msg.ValidatorDstAddress),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(stakingTypes.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, stakingTypes.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress),
		),
	})
	return &types.MsgBeginRedelegateResponse{
		CompletionTime: completionTime,
	}, nil
}

func ParseBech32ValConsPubkey(validatorInfoPubKeyBase64 string) (cryptotypes.PubKey, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainMsgServer)
	validatorInfoPubKeyBytes, err := base64.StdEncoding.DecodeString(validatorInfoPubKeyBase64)
	if err != nil {
		log.WithError(err).Error("decode validator info pub key failed")
		return nil, err
	}
	pbk := ed25519.PubKey(validatorInfoPubKeyBytes)
	pubkey, err := cryptocodec.FromTmPubKeyInterface(pbk)
	if err != nil {
		log.WithError(err).Error("from tm pubkey interface failed")
		return nil, err
	}
	return pubkey, nil
}
