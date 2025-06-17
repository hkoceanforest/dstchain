package keeper

import (
	"fmt"
	"freemasonry.cc/blockchain/core"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingKeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/tendermint/tendermint/libs/log"
	"time"

	"freemasonry.cc/blockchain/x/gateway/types"
)

type Keeper struct {
	storeKey   storetypes.StoreKey
	cdc        codec.BinaryCodec
	paramstore paramtypes.Subspace

	stakingKeeper *stakingKeeper.Keeper
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	hooks         types.CommonHooks
}

func NewKeeper(
	storeKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
	ps paramtypes.Subspace,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	stakingKeeper *stakingKeeper.Keeper,
) Keeper {
	
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeKey:      storeKey,
		cdc:           cdc,
		paramstore:    ps,
		accountKeeper: ak,
		bankKeeper:    bk,
		stakingKeeper: stakingKeeper,
		hooks:         nil,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) KVHelper(ctx sdk.Context) StoreHelper {
	store := ctx.KVStore(k.storeKey)
	return StoreHelper{
		store,
	}
}

func (k *Keeper) SetHooks(sh types.CommonHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set common hooks twice")
	}
	k.hooks = sh
	return k
}

func (k Keeper) GetGatewayUpload(ctx sdk.Context, validatorAddress string) []byte {
	
	kvStore := ctx.KVStore(k.storeKey)
	key := []byte(types.GatewayUploadKey + validatorAddress)

	if kvStore.Has(key) {
		res := kvStore.Get(key)
		return res
	}

	return nil
}


func (k Keeper) SetGatewayDelegateLastHeight(ctx sdk.Context, delegateAddress, validatorAddress string) error {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainCommKeeper)
	store := k.KVHelper(ctx)
	key := types.DelegateLastTimeKey + delegateAddress + "_" + validatorAddress
	err := store.Set(key, ctx.BlockHeight())
	if err != nil {
		log.WithError(err).Error("SetGatewayDelegateLastTime")
		return err
	}
	return nil
}


func (k Keeper) GetGatewayDelegateLastHeight(ctx sdk.Context, delegateAddress, validatorAddress string) (int64, error) {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainCommKeeper)
	store := k.KVHelper(ctx)
	key := types.DelegateLastTimeKey + delegateAddress + "_" + validatorAddress
	var lastHeight int64
	if !store.Has(key) {
		return lastHeight, nil
	}
	err := store.GetUnmarshal(key, &lastHeight)
	if err != nil {
		log.WithError(err).Error("GetGatewayDelegateLastTime")
		return 0, err
	}
	return lastHeight, nil
}




func (k Keeper) RedeemCheck(ctx sdk.Context, params types.Params) error {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainCommKeeper)
	if ctx.BlockHeight()%params.IndexNumHeight == 0 {
		redeemMap, err := k.GetGatewayRedeemNum(ctx)
		if err != nil {
			log.WithError(err).Error("GetGatewayRedeemNum Err")
			return err
		}
		var numArray []types.GatewayNumIndex
		for _, val := range redeemMap {
			
			if val.Status == 1 && val.Validity <= ctx.BlockHeight() {
				val.GatewayAddress = ""
				val.Status = 2
				val.Validity = 0
				val.IsFirst = false
				numArray = append(numArray, val)
			}
		}
		
		err = k.SetGatewayNum(ctx, numArray)
		if err != nil {
			log.WithError(err).Error("SetGatewayNum Err")
			return err
		}
		
		err = k.GatewayRedeemNumFilter(ctx, numArray)
		if err != nil {
			log.WithError(err).Error("GatewayRedeemNumFilter Err")
			return err
		}
	}

	return nil
}


func (k Keeper) createValidator(ctx sdk.Context, delegatorAddress sdk.AccAddress, validatorAddress sdk.ValAddress, msg types.MsgCreateSmartValidator, delegation sdk.Coin) error {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainCommKeeper)
	pk, err := ParseBech32ValConsPubkey(msg.PubKey)
	if err != nil {
		log.WithError(err).Error("ParseBech32ValConsPubkey Err")
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "Expecting cryptotypes.PubKey, got %T", pk)
	}

	if _, found := k.stakingKeeper.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(pk)); found {
		return stakingTypes.ErrValidatorPubKeyExists
	}
	bondDenom := k.stakingKeeper.BondDenom(ctx)
	if delegation.Denom != bondDenom {
		log.WithError(sdkerrors.ErrInvalidRequest).Error("createValidator denom err")
		return sdkerrors.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s, expected %s", delegation.Denom, bondDenom,
		)
	}

	if _, err := msg.Description.EnsureLength(); err != nil {
		log.WithError(err).Error("EnsureLength Err")
		return err
	}
	validator, err := stakingTypes.NewValidator(validatorAddress, pk, msg.Description)
	if err != nil {
		log.WithError(err).Error("NewValidator Err")
		return err
	}
	commission := stakingTypes.NewCommissionWithTime(
		msg.Commission.Rate, msg.Commission.MaxRate,
		msg.Commission.MaxChangeRate, ctx.BlockHeader().Time,
	)
	validator, err = validator.SetInitialCommission(commission)
	
	validator.MinSelfDelegation = msg.MinSelfDelegation
	k.stakingKeeper.SetValidator(ctx, validator)
	k.stakingKeeper.SetValidatorByConsAddr(ctx, validator)
	k.stakingKeeper.SetNewValidatorByPowerIndex(ctx, validator)
	k.stakingKeeper.AfterValidatorCreated(ctx, validator.GetOperator())
	
	_, err = k.stakingKeeper.Delegate(ctx, delegatorAddress, delegation.Amount, stakingTypes.Unbonded, validator, true)
	if err != nil {
		log.WithError(err).Error("Delegate Err")
		return err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			stakingTypes.EventTypeCreateValidator,
			sdk.NewAttribute(stakingTypes.AttributeKeyValidator, validator.GetOperator().String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, delegation.String()),
		),
	})
	return nil
}

func (k Keeper) delegate(ctx sdk.Context, delegatorAddress sdk.AccAddress, validatorAddress sdk.ValAddress, validator stakingTypes.Validator, coin sdk.Coin) error {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainCommKeeper)
	bondDenom := k.stakingKeeper.BondDenom(ctx)
	if coin.Denom != bondDenom {
		log.WithError(sdkerrors.ErrInvalidRequest).Error("delegate denom err")
		return sdkerrors.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s, expected %s", coin.Denom, bondDenom,
		)
	}
	
	newShares, err := k.stakingKeeper.Delegate(ctx, delegatorAddress, coin.Amount, stakingTypes.Unbonded, validator, true)
	if err != nil {
		log.WithError(err).Error("Delegate Err")
		return err
	}
	
	err = k.SetGatewayDelegateLastHeight(ctx, delegatorAddress.String(), validatorAddress.String())
	if err != nil {
		log.WithError(err).Error("SetGatewayDelegateLastTime Err")
		return err
	}
	coins := k.bankKeeper.GetAllBalances(ctx, delegatorAddress)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			stakingTypes.EventTypeDelegate,
			sdk.NewAttribute(stakingTypes.AttributeKeyValidator, validatorAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, coin.String()),
			sdk.NewAttribute(stakingTypes.AttributeKeyNewShares, newShares.String()),    
			sdk.NewAttribute(stakingTypes.AttributeKeyDelegatorBalance, coins.String()), 
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, stakingTypes.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, delegatorAddress.String()),
		),
	})
	return nil
}

func (k Keeper) Undelegate(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, validator stakingTypes.Validator, shares sdk.Dec) (time.Time, sdk.Int, error) {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainCommKeeper)
	if k.stakingKeeper.HasMaxUnbondingDelegationEntries(ctx, delAddr, valAddr) {
		return time.Time{}, sdk.ZeroInt(), stakingTypes.ErrMaxUnbondingDelegationEntries
	}
	returnAmount, err := k.stakingKeeper.Unbond(ctx, delAddr, valAddr, shares)
	if err != nil {
		log.WithError(err).Error("Unbond Err")
		return time.Time{}, sdk.ZeroInt(), err
	}
	
	if validator.GetOperator().String() == sdk.ValAddress(delAddr).String() {
		params := k.GetParams(ctx)
		
		lastHeight, err := k.GetGatewayDelegateLastHeight(ctx, delAddr.String(), valAddr.String())
		if err != nil {
			log.WithError(err).Error("GetGatewayDelegateLastTime Err")
			return time.Time{}, sdk.ZeroInt(), err
		}
		diff := ctx.BlockHeight() - lastHeight
		if diff < params.RedeemFeeHeight { 
			fees := sdk.NewDecFromInt(returnAmount).Mul(params.RedeemFee)
			returnAmount = returnAmount.Sub(fees.RoundInt())
			coin := sdk.NewCoin(sdk.DefaultBondDenom, fees.RoundInt())
			err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, stakingTypes.BondedPoolName, authtypes.FeeCollectorName, sdk.NewCoins(coin))
			if err != nil {
				log.WithError(err).Error("SendCoinsFromModuleToModule Err")
				panic(err)
			}
		}
	}
	
	if validator.IsBonded() {
		k.bondedTokensToNotBonded(ctx, returnAmount)
	}

	completionTime := ctx.BlockHeader().Time.Add(k.stakingKeeper.UnbondingTime(ctx))
	ubd := k.stakingKeeper.SetUnbondingDelegationEntry(ctx, delAddr, valAddr, ctx.BlockHeight(), completionTime, returnAmount)
	k.stakingKeeper.InsertUBDQueue(ctx, ubd, completionTime)

	return completionTime, returnAmount, nil
}

func (k Keeper) bondedTokensToNotBonded(ctx sdk.Context, tokens sdk.Int) {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainCommKeeper)
	coins := sdk.NewCoins(sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), tokens))
	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, stakingTypes.BondedPoolName, stakingTypes.NotBondedPoolName, coins); err != nil {
		log.WithError(err).Error("SendCoinsFromModuleToModule Err")
		panic(err)
	}
}
