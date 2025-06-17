package client

import (
	"context"
	"cosmossdk.io/math"
	"encoding/json"
	"errors"
	"fmt"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/core/chainnet"
	"freemasonry.cc/blockchain/util"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/client/common"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"sort"
	"strings"
	"time"
)

type DposClient struct {
	TxClient  *TxClient
	logPrefix string
}

func (this *DposClient) RegisterValidator(bech32DelegatorAddr, bech32ValidatorAddr string, bech32ValidatorPubkey cryptotypes.PubKey, selfDelegation sdk.Coin, desc stakingTypes.Description, commission stakingTypes.CommissionRates, minSelfDelegation sdk.Int, privateKey string, fee float64) (resp *core.BaseResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	_, err = sdk.AccAddressFromBech32(bech32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		err = core.ParseAccountError
		return
	}
	validatorAddr, err := sdk.ValAddressFromBech32(bech32ValidatorAddr)
	if err != nil {
		log.WithError(err).Error("ValAddressFromBech32")
		err = core.ParseAccountError
		return
	}

	
	err = commission.Validate()
	if err != nil {
		log.WithError(err).Error("commission.Validate")
		return
	}
	msg, err := stakingTypes.NewMsgCreateValidator(validatorAddr, bech32ValidatorPubkey, selfDelegation, desc, commission, minSelfDelegation)
	if err != nil {
		log.WithError(err).Error("NewMsgCreateValidator")
		return
	}
	_, resp, err = this.TxClient.SignAndSendMsg(bech32DelegatorAddr, privateKey, core.NewLedgerFee(fee), "", msg)
	if err != nil {
		return
	}
	if resp.Status == 1 {
		return resp, nil 
	} else {
		return resp, errors.New(resp.Info) 
	}
}


func (this *DposClient) EditorValidator(bech32ValidatorAccAddr string, desc stakingTypes.Description, newRate *sdk.Dec, minSelfDelegation *sdk.Int, privateKey string) (resp *core.BaseResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	accAddr, err := sdk.AccAddressFromBech32(bech32ValidatorAccAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		return
	}
	validatorAddress := sdk.ValAddress(accAddr).String()
	
	validatorInfor, err := this.FindValidatorByValAddress(validatorAddress)
	if err != nil {
		err = errors.New(core.QueryChainInforError)
		return
	}
	
	if validatorInfor.GetOperator().Empty() {
		err = stakingTypes.ErrNoValidatorFound
		return
	}
	if minSelfDelegation != nil && !minSelfDelegation.Equal(validatorInfor.MinSelfDelegation) {
		if !minSelfDelegation.GT(validatorInfor.MinSelfDelegation) {
			return nil, stakingTypes.ErrMinSelfDelegationDecreased
		}
		if minSelfDelegation.GT(validatorInfor.Tokens) {
			return nil, stakingTypes.ErrSelfDelegationBelowMinimum
		}
	} else {
		minSelfDelegation = nil
	}
	
	
	if time.Now().Sub(validatorInfor.Commission.UpdateTime).Hours() < 24 {
		err = errors.New(core.ValidatorInfoError)
		return
	}
	_, err = desc.EnsureLength()
	if err != nil {
		log.WithError(err).Error("EnsureLength")
		err = errors.New(core.ValidatorDescriptionError)
		return
	}
	msg := stakingTypes.NewMsgEditValidator(validatorInfor.GetOperator(), desc, newRate, minSelfDelegation)
	if err != nil {
		log.WithError(err).Error("NewMsgEditValidator")
		return
	}
	_, resp, err = this.TxClient.SignAndSendMsg(bech32ValidatorAccAddr, privateKey, core.NewLedgerFee(0), "", msg)
	if err != nil {
		return
	}
	if resp.Status == 1 {
		return resp, nil 
	} else {
		return resp, errors.New(resp.Info) 
	}
}


func (this *DposClient) FindValidatorByValAddress(bech32ValidatorAddr string) (validator *stakingTypes.Validator, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	validatorAddr, err := sdk.ValAddressFromBech32(bech32ValidatorAddr)
	if err != nil {
		log.WithError(err).Error("ValAddressFromBech32")
		err = core.ParseAccountError
		return
	}
	params := stakingTypes.QueryValidatorParams{ValidatorAddr: validatorAddr}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, err
	}
	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/staking/"+stakingTypes.QueryValidator, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return nil, err
	}
	validator = &stakingTypes.Validator{}

	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, validator)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
	}
	return
}


func (this *DposClient) UnjailValidator(bech32DelegatorAddr, bech32ValidatorAddr, privateKey string) (resp *core.BaseResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	validatorAddr, err := sdk.ValAddressFromBech32(bech32ValidatorAddr)
	if err != nil {
		log.WithError(err).Error("ValAddressFromBech32")
		return
	}
	
	validatorInfo, err := this.FindValidatorByValAddress(bech32ValidatorAddr)
	if err != nil {
		log.WithError(err).Error("FindValidatorByValAddress")
		return
	}
	if !validatorInfo.Jailed {
		return
	}

	accAddr, err := sdk.AccAddressFromBech32(bech32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		return
	}
	OperatorAddress := sdk.ValAddress(accAddr).String()
	if validatorInfo.OperatorAddress != OperatorAddress {
		log.Error("validatorInfo.OperatorAddress:", validatorInfo.OperatorAddress, "|OperatorAddress:", OperatorAddress)
		return nil, core.ErrUnjailOperatorAddress
	}
	
	delegatorResponse, _, err := this.FindDelegation(bech32DelegatorAddr, bech32ValidatorAddr)
	if err != nil {
		log.WithError(err).Error("this.FindDelegatio", bech32DelegatorAddr, bech32ValidatorAddr)
		return
	}
	
	if delegatorResponse.Delegation.Shares.IsZero() {
		log.Error("delegatorResponse.Delegation.Shares.IsZero")
		return nil, core.ErrDelegationZero
	}

	tokens := validatorInfo.TokensFromShares(delegatorResponse.Delegation.Shares).TruncateInt()
	if tokens.LT(validatorInfo.MinSelfDelegation) {
		log.Error("tokens.LT(validatorInfo.MinSelfDelegation)")
		return nil, core.ErrDelegateAmountLtMinSelfDelegation
	}

	
	msg := slashingTypes.NewMsgUnjail(validatorAddr)
	fee := legacytx.NewStdFee(flags.DefaultGasLimit, sdk.NewCoins(sdk.NewCoin(core.BaseDenom, sdk.NewInt(0))))
	_, resp, err = this.TxClient.SignAndSendMsg(bech32DelegatorAddr, privateKey, fee, "", msg)
	if err != nil {
		log.WithError(err).Error("this.TxClient.SignAndSendMsg")
		return
	}
	if resp.Status == 1 {
		return resp, nil 
	} else {
		return resp, errors.New(resp.Info) 
	}
}


func (this *DposClient) FindDelegation(delegatorAddr, validatorAddr string) (delegation *stakingTypes.DelegationResponse, notFound bool, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	notFound = false
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	params := stakingTypes.QueryDelegatorValidatorRequest{DelegatorAddr: delegatorAddr, ValidatorAddr: validatorAddr}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, notFound, err
	}
	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/staking/"+stakingTypes.QueryDelegation, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		if strings.Contains(err.Error(), stakingTypes.ErrNoDelegation.Error()) {
			notFound = true
		}
		return nil, notFound, err
	}
	delegation = &stakingTypes.DelegationResponse{}
	err = util.Json.Unmarshal(resBytes, delegation)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
	}
	return
}

type ValidatorInfo struct {
	stakingTypes.Validator
	StartTime           int64  `json:"start_time"`
	UnbondingTimeFormat string `json:"unbonding_time_format"`
}


func (this *DposClient) QueryValidators(page, limit int, status string) ([]ValidatorInfo, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	if status == "" {
		status = stakingTypes.BondStatusBonded
	}

	params := stakingTypes.NewQueryValidatorsParams(page, limit, status)

	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("clientCtx.LegacyAmino.MarshalJSON")
		return nil, err
	}

	route := fmt.Sprintf("custom/%s/%s", stakingTypes.QuerierRoute, stakingTypes.QueryValidators)

	res, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, route, bz)
	if err != nil {
		log.WithError(err).Error("clientCtx.QueryWithData")
		return nil, err
	}

	validatorsResp := make(stakingTypes.Validators, 0)
	err = clientCtx.LegacyAmino.UnmarshalJSON(res, &validatorsResp)
	if err != nil {
		return nil, err
	}

	validatorInfos := make([]ValidatorInfo, 0)

	for _, val := range validatorsResp {

		valConsAddr, err := val.GetConsAddr()
		if err != nil {
			log.WithError(err).Error("GetConsAddr Err:" + err.Error())
			return nil, err
		}

		
		node, err := clientCtx.GetNode()
		if err != nil {
			log.WithError(err).Error("clientCtx.GetNode")
			return nil, err
		}
		slashingParams := slashingTypes.QuerySigningInfoRequest{valConsAddr.String()}
		bz, err := clientCtx.LegacyAmino.MarshalJSON(slashingParams)
		if err != nil {
			log.WithError(err).Error("MarshalJSON")
			return nil, err
		}

		resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/slashing/"+slashingTypes.QuerySigningInfo, bz)
		if err != nil {
			log.WithError(err).Error("QueryWithData")
			return nil, err
		}
		info := slashingTypes.ValidatorSigningInfo{}
		err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &info)
		if err != nil {
			log.WithError(err).Error("clientCtx.LegacyAmino.UnmarshalJSON(slashingTypes.ValidatorSigningInfo)")
			return nil, err
		}
		nodeStatus, err := node.Status(context.Background())
		if err != nil {
			log.WithError(err).Error("node.Status")
			return nil, err
		}
		startHeight := nodeStatus.SyncInfo.LatestBlockHeight - info.IndexOffset
		if startHeight == 0 {
			startHeight = 1
		}
		blockInfo, err := node.Block(context.Background(), &startHeight)
		if err != nil {
			log.WithError(err).Error("node.Block")
			return nil, err
		}

		
		unbondingTimeFormat := ""
		if val.UnbondingTime.Unix() != 0 {
			unbondingTimeFormat = val.UnbondingTime.Format("2006-01-02 15:04")
		}

		ValidatorInfoEdit := ValidatorInfo{
			val,
			time.Now().Unix() - blockInfo.Block.Time.Unix(),
			unbondingTimeFormat,
		}

		ValidatorInfoEdit.ConsensusPubkey = nil
		validatorInfos = append(validatorInfos, ValidatorInfoEdit)

	}

	return validatorInfos, nil

}

type DecCoinsString struct {
	Coins []struct {
		Denom  string
		Amount string
	} `json:"coins"`
}


func (this *DposClient) QueryValCanWithdraw(accAddr string) (res distributionTypes.ValidatorAccumulatedCommission, err error) {
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	

	accFromAddress, err := sdk.AccAddressFromBech32(accAddr)
	if err != nil {
		return res, err
	}
	valAddr := sdk.ValAddress(accFromAddress)

	
	bz, err := common.QueryValidatorCommission(clientCtx, valAddr)
	if err != nil {
		return res, err
	}

	var commission distributionTypes.ValidatorAccumulatedCommission
	clientCtx.LegacyAmino.UnmarshalJSON(bz, &commission)

	return commission, nil
}


func (this *DposClient) Delegation(bech32DelegatorAddr, bech32ValidatorAddr string, amount sdk.Coin, privateKey string) (resp *core.BaseResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	delegatorAddr, err := sdk.AccAddressFromBech32(bech32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		err = core.ParseAccountError
		return
	}
	validatorAddr, err := sdk.ValAddressFromBech32(bech32ValidatorAddr)
	if err != nil {
		log.WithError(err).Error("ValAddressFromBech32")
		err = core.ParseAccountError
		return
	}
	msg := stakingTypes.NewMsgDelegate(delegatorAddr, validatorAddr, amount)
	fee := core.NewLedgerFee(0)
	_, resp, err = this.TxClient.SignAndSendMsg(msg.DelegatorAddress, privateKey, fee, "", msg)
	if err != nil {
		return
	}
	if resp.Status == 1 {
		return resp, nil 
	} else {
		return resp, errors.New(resp.Info) 
	}
}

func (this *DposClient) UnbondDelegation(bech32DelegatorAddr string, bech32ValidatorAddr string, amount sdk.Coin, privateKey string) (resp *core.BaseResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	delegatorAddr, err := sdk.AccAddressFromBech32(bech32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		err = core.ParseAccountError
		return
	}
	validatorAddr, err := sdk.ValAddressFromBech32(bech32ValidatorAddr)
	if err != nil {
		log.WithError(err).Error("ValAddressFromBech32")
		err = core.ParseAccountError
		return
	}
	msg := stakingTypes.NewMsgUndelegate(delegatorAddr, validatorAddr, amount)
	fee := core.NewLedgerFee(0)
	_, resp, err = this.TxClient.SignAndSendMsg(msg.DelegatorAddress, privateKey, fee, "", msg)
	if err != nil {
		return
	}
	if resp.Status == 1 {
		return resp, nil 
	} else {
		return resp, errors.New(resp.Info) 
	}
}


func (this *DposClient) UnbondDelegationAll(bech32DelegatorAddr, privateKey string) (resp *core.BaseResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	delegatorAddr, err := sdk.AccAddressFromBech32(bech32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		err = core.ParseAccountError
		return
	}
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	params := stakingTypes.NewQueryDelegatorParams(delegatorAddr)
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, err
	}
	
	resBytes, _, err := clientCtx.QueryWithData("custom/staking/"+stakingTypes.QueryDelegatorDelegations, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return
	}
	delegationResp := stakingTypes.DelegationResponses{}
	err = util.Json.Unmarshal(resBytes, &delegationResp)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
		return
	}
	if len(delegationResp) <= 0 {
		err = errors.New("delegation does not exist")
		return
	}
	msgs := []sdk.Msg{}
	for _, delegation := range delegationResp {
		valAddr, err := sdk.ValAddressFromBech32(delegation.Delegation.ValidatorAddress)
		if err != nil {
			continue
		}
		msg := stakingTypes.NewMsgUndelegate(delegatorAddr, valAddr, delegation.Balance)
		msgs = append(msgs, msg)
	}

	fee := core.NewLedgerFee(0)
	_, resp, err = this.TxClient.SignAndSendMsg(bech32DelegatorAddr, privateKey, fee, "", msgs...)
	if err != nil {
		return
	}
	if resp.Status == 1 {
		return resp, nil 
	} else {
		return resp, errors.New(resp.Info) 
	}
}


func (this *DposClient) DrawCommissionDelegationRewards(bech32DelegatorAddr, bech32ValidatorAddr, privateKey string) (resp *core.BaseResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	delegatorAddr, err := sdk.AccAddressFromBech32(bech32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		err = core.ParseAccountError
		return
	}
	validatorAddr, err := sdk.ValAddressFromBech32(bech32ValidatorAddr)
	if err != nil {
		log.WithError(err).Error("ValAddressFromBech32")
		err = core.ParseAccountError
		return
	}
	bech32DelegatorValidatorAddr := sdk.ValAddress(delegatorAddr).String()

	
	delegationReward, validatorReward, err := this.RewardsPreview(bech32DelegatorAddr, bech32ValidatorAddr)
	if err != nil {
		return
	}

	

	if validatorReward.IsZero() && delegationReward.IsZero() {
		err = core.ErrRewardReceive 
		return
	}
	msgs := []sdk.Msg{}

	if !delegationReward.IsZero() { 
		msg1 := distributionTypes.NewMsgWithdrawDelegatorReward(delegatorAddr, validatorAddr)
		msgs = append(msgs, msg1)
	}
	
	if bech32DelegatorValidatorAddr == bech32ValidatorAddr && !validatorReward.IsZero() { 
		
		msg2 := distributionTypes.NewMsgWithdrawValidatorCommission(validatorAddr)
		msgs = append(msgs, msg2) 
	}

	fee := core.NewLedgerFee(0)
	_, resp, err = this.TxClient.SignAndSendMsg(bech32DelegatorAddr, privateKey, fee, "", msgs...)
	if err != nil {
		return
	}
	if resp.Status == 1 {
		return resp, nil 
	} else {
		return resp, errors.New(resp.Info) 
	}
}


func (this *DposClient) RewardsPreview(bech32DelegatorAddr string, bech32ValidatorAddr string) (delegatorReward sdk.Coin, validatorReward sdk.Coin, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	delegatorReward = sdk.NewCoin(core.GovDenom, sdk.NewInt(0))
	validatorReward = sdk.NewCoin(core.GovDenom, sdk.NewInt(0))
	delegatorAddr, err := sdk.AccAddressFromBech32(bech32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		return delegatorReward, validatorReward, core.ParseAccountError
	}
	validatorAddr, err := sdk.ValAddressFromBech32(bech32ValidatorAddr)
	if err != nil {
		log.WithError(err).Error("ValAddressFromBech32")
		return delegatorReward, validatorReward, core.ParseAccountError
	}
	params := distributionTypes.NewQueryDelegationRewardsParams(delegatorAddr, validatorAddr)
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return delegatorReward, validatorReward, errors.New(core.QueryChainInforError)
	}
	
	resBytes, _, err := clientCtx.QueryWithData("custom/distribution/"+distributionTypes.QueryDelegationRewards, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData error1 | ", err.Error())
		return delegatorReward, validatorReward, err
	}

	delegatorRewardDecCoins := sdk.DecCoins{} 

	err = util.Json.Unmarshal(resBytes, &delegatorRewardDecCoins)
	if err != nil {
		log.WithError(err).Error("Unmarshal error1 | ", err.Error())
		return delegatorReward, validatorReward, errors.New(core.QueryChainInforError)
	}

	if len(delegatorRewardDecCoins) != 0 {
		delegatorReward, _ = delegatorRewardDecCoins[0].TruncateDecimal() 
	}
	bech32DelegatorValidatorAddr := sdk.ValAddress(delegatorAddr).String()
	
	if bech32DelegatorValidatorAddr == bech32ValidatorAddr {
		
		resBytes, _, err = clientCtx.QueryWithData("custom/distribution/"+distributionTypes.QueryValidatorCommission, bz)
		if err != nil {
			log.WithError(err).Error("QueryWithData error2 | ", err.Error())
			return delegatorReward, validatorReward, errors.New(core.QueryChainInforError)
		}
		validatorCommDecCoins := distributionTypes.ValidatorAccumulatedCommission{} 
		err = util.Json.Unmarshal(resBytes, &validatorCommDecCoins)
		if err != nil {
			log.WithError(err).Error("Unmarshal error2 | ", err.Error())
			return delegatorReward, validatorReward, errors.New(core.QueryChainInforError)
		}
		if len(validatorCommDecCoins.Commission) > 0 {
			validatorReward, _ = validatorCommDecCoins.Commission[0].TruncateDecimal() 
		}
	}

	return delegatorReward, validatorReward, err
}


func (this *DposClient) DrawCommissionDelegationRewardsAll(bech32DelegatorAddr, privateKey string) (resp *core.BaseResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	delegatorAddr, err := sdk.AccAddressFromBech32(bech32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		err = core.ParseAccountError
		return
	}
	totalRewards, ValAddressArr, err := this.RewardsPreviewAll(bech32DelegatorAddr)
	if err != nil {
		err = core.ParseAccountError
		return
	}
	if totalRewards.Total.IsZero() && len(ValAddressArr.ValidatorCommissions) <= 0 {
		err = errors.New("There is no reward to receive") 
		return
	}
	msgs := []sdk.Msg{}

	if !totalRewards.Total.IsZero() { 
		for _, reward := range totalRewards.Rewards {
			
			if reward.Reward.IsZero() {
				continue
			}
			if reward.Reward[0].Amount.LT(sdk.NewDec(1)) {
				continue
			}
			validatorAddr, err := sdk.ValAddressFromBech32(reward.ValidatorAddress)
			if err != nil {
				continue
			}
			msg1 := distributionTypes.NewMsgWithdrawDelegatorReward(delegatorAddr, validatorAddr)
			msgs = append(msgs, msg1)
		}
	}
	
	if len(ValAddressArr.ValidatorCommissions) > 0 { 
		
		for _, valAddress := range ValAddressArr.ValidatorCommissions {
			msg2 := distributionTypes.NewMsgWithdrawValidatorCommission(valAddress.ValidatorAddress)
			msgs = append(msgs, msg2) 
		}
	}

	fee := core.NewLedgerFee(0)
	_, resp, err = this.TxClient.SignAndSendMsg(bech32DelegatorAddr, privateKey, fee, "", msgs...)
	if err != nil {
		return
	}
	if resp.Status == 1 {
		return resp, nil 
	} else {
		return resp, errors.New(resp.Info) 
	}
}


func (this *DposClient) RewardsPreviewAll(bech32DelegatorAddr string) (distributionTypes.QueryDelegatorTotalRewardsResponse, stakingTypes.ValidatorCommissionResp, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	totalRewards := distributionTypes.QueryDelegatorTotalRewardsResponse{}
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	var ValAddressArr stakingTypes.ValidatorCommissionResp
	var valAddrcommission stakingTypes.ValidatorCommission
	ValAddressArr.Total = sdk.DecCoins{}
	delegatorAddr, err := sdk.AccAddressFromBech32(bech32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		err = core.ParseAccountError
		return totalRewards, ValAddressArr, nil
	}
	params := distributionTypes.NewQueryDelegatorParams(delegatorAddr)
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON error1 | ", err.Error())
		return totalRewards, ValAddressArr, nil
	}
	
	resBytes, _, err := clientCtx.QueryWithData("custom/distribution/"+distributionTypes.QueryDelegatorTotalRewards, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData error1 | ", err.Error())
		return totalRewards, ValAddressArr, nil
	}

	err = util.Json.Unmarshal(resBytes, &totalRewards)
	if err != nil {
		log.WithError(err).Error("Unmarshal error1 | ", err.Error())
		return totalRewards, ValAddressArr, nil
	}
	for _, reward := range totalRewards.Rewards {
		validatorAddr, err := sdk.ValAddressFromBech32(reward.ValidatorAddress)
		if err != nil {
			log.WithError(err).Error("ValAddressFromBech32 error2 | ", err.Error())
			continue
		}
		params := distributionTypes.NewQueryDelegationRewardsParams(delegatorAddr, validatorAddr)
		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			log.WithError(err).Error("MarshalJSON error2 | ", err.Error())
			continue
		}
		bech32DelegatorValidatorAddr := sdk.ValAddress(delegatorAddr).String()
		
		if bech32DelegatorValidatorAddr == validatorAddr.String() {
			
			resBytes, _, err = clientCtx.QueryWithData("custom/distribution/"+distributionTypes.QueryValidatorCommission, bz)
			if err != nil {
				log.WithError(err).Error("QueryWithData error2 | ", err.Error())
				continue
			}
			validatorComm := distributionTypes.ValidatorAccumulatedCommission{} 
			err = util.Json.Unmarshal(resBytes, &validatorComm)
			if err != nil {
				log.WithError(err).Error("Unmarshal error2 | ", err.Error())
				continue
			}
			if validatorComm.Commission.IsZero() {
				continue
			}
			valAddrcommission.ValidatorAddress = validatorAddr
			valAddrcommission.Reward = validatorComm.Commission
			for _, coin := range validatorComm.Commission {
				ValAddressArr.Total = ValAddressArr.Total.Add(coin)
			}
			ValAddressArr.ValidatorCommissions = append(ValAddressArr.ValidatorCommissions, valAddrcommission)
		}

	}
	return totalRewards, ValAddressArr, nil
}

type DposInfo struct {
	AllStaking         math.Int                 `json:"all_staking"`
	AllSupply          sdk.Coin                 `json:"all_supply"`
	StakingParams      stakingTypes.Params      `json:"staking_params"`
	DistributionParams distributionTypes.Params `json:"distribution_params"`
	SlashingParams     slashingTypes.Params     `json:"slashing_params"`
	MintInflation      string                   `json:"mint_inflation"`
}

func (this *DposClient) GetDposParams() (DposInfo, error) {

	var res DposInfo
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	
	resBytes, _, err := clientCtx.QueryWithData("custom/staking/"+stakingTypes.QueryPool, nil)
	if err != nil {
		return DposInfo{}, err
	}
	stakingPool := stakingTypes.Pool{}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &stakingPool)
	if err != nil {
		return DposInfo{}, err
	}

	res.AllStaking = stakingPool.BondedTokens

	
	supplyParams := bankTypes.QuerySupplyOfParams{
		Denom: "nxn",
	}
	supplyParamsData, err := clientCtx.LegacyAmino.MarshalJSON(supplyParams)

	resBytes, _, err = clientCtx.QueryWithData("custom/bank/"+bankTypes.QuerySupplyOf, supplyParamsData)
	if err != nil {
		return DposInfo{}, err
	}
	supplyCoin := sdk.Coin{}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &supplyCoin)
	if err != nil {
		return DposInfo{}, err
	}

	res.AllSupply = supplyCoin

	
	resBytes, _, err = clientCtx.QueryWithData("custom/staking/"+stakingTypes.QueryParameters, nil)
	if err != nil {
		return DposInfo{}, err
	}
	stakingParams := stakingTypes.Params{}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &stakingParams)
	if err != nil {
		return DposInfo{}, err
	}

	res.StakingParams = stakingParams

	
	resBytes, _, err = clientCtx.QueryWithData("custom/distribution/"+distributionTypes.QueryParams, nil)
	if err != nil {
		return DposInfo{}, err
	}
	distriParams := distributionTypes.Params{}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &distriParams)
	if err != nil {
		return DposInfo{}, err
	}

	res.DistributionParams = distriParams

	
	resBytes, _, err = clientCtx.QueryWithData("custom/slashing/"+slashingTypes.QueryParameters, nil)
	if err != nil {
		return DposInfo{}, err
	}
	slashingParams := slashingTypes.Params{}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &slashingParams)
	if err != nil {
		return DposInfo{}, err
	}

	res.SlashingParams = slashingParams

	
	resBytes, _, err = clientCtx.QueryWithData("custom/mint/"+minttypes.QueryInflation, nil)
	if err != nil {
		return DposInfo{}, err
	}
	mintInflation := sdk.ZeroDec()
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &mintInflation)
	if err != nil {
		return DposInfo{}, err
	}

	res.MintInflation = core.RemoveStringLastZero(mintInflation.String())

	return res, nil
}

type DposValidatorInfo struct {
	
	ValidatorAddr string `json:"validator_addr"`
	
	ValidatorName string `json:"validator_name"`
	
	ValidatorDelegateAmount math.Int `json:"validator_delegate_amount"`
	
	ValidatorOnlineRatio string `json:"validator_online_ratio"`
	
	ValidatorCommision sdk.Dec `json:"commision"`
	
	PersonDelegateAmount math.Int `json:"persion_delegate_amount"`
}

type DposValidatorList struct {
	
	ValidatorAmount int64 `json:"validator_amount"`
	
	ValidatorAmountMax uint32 `json:"validator_amount_max"`
	
	ValidatorList []DposValidatorInfo `json:"validator_list"`
}


func (this *DposClient) ValidatorList(fromAddr string, isDelegateFlag, sortType string, page, limit int64, valName string) (resp DposValidatorList, err error) {
	resp.ValidatorList = make([]DposValidatorInfo, 0)
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	accFromAddr, err := sdk.AccAddressFromBech32(fromAddr)
	if err != nil {
		return
	}

	
	resBytes, _, err := clientCtx.QueryWithData("custom/staking/"+stakingTypes.QueryParameters, nil)
	if err != nil {
		return
	}
	stakingParams := stakingTypes.Params{}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &stakingParams)
	if err != nil {
		return
	}

	resp.ValidatorAmountMax = stakingParams.MaxValidators

	
	requestValidatorListParams := stakingTypes.QueryValidatorsParams{
		Page:   1,
		Limit:  int(stakingParams.MaxValidators),
		Status: "BOND_STATUS_BONDED",
	}
	requestValidatorListParamsData, err := clientCtx.LegacyAmino.MarshalJSON(requestValidatorListParams)
	if err != nil {
		return
	}

	resBytes, _, err = clientCtx.QueryWithData("custom/staking/"+stakingTypes.QueryValidators, requestValidatorListParamsData)
	if err != nil {
		return
	}
	valList := []stakingTypes.Validator{}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &valList)
	if err != nil {
		return
	}

	
	requestDelegationsParams := stakingTypes.QueryDelegatorParams{
		DelegatorAddr: accFromAddr,
	}
	requestDelegationsParamsData, err := clientCtx.LegacyAmino.MarshalJSON(requestDelegationsParams)
	if err != nil {
		return
	}

	resBytes, _, err = clientCtx.QueryWithData("custom/staking/"+stakingTypes.QueryDelegatorDelegations, requestDelegationsParamsData)
	if err != nil {
		return
	}
	delegations := stakingTypes.DelegationResponses{}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &delegations)
	if err != nil {
		return
	}

	delegationsMap := make(map[string]sdk.Coin)

	for _, del := range delegations {
		delegationsMap[del.Delegation.ValidatorAddress] = del.Balance
	}

	
	var valInfos []DposValidatorInfo
	var isDelegate bool
	for _, validator := range valList {
		delVal := sdk.ZeroInt()
		delegate, ok := delegationsMap[validator.OperatorAddress]
		if ok == true {
			delVal = delegate.Amount
		}

		if valName != "" {
			if !strings.Contains(validator.Description.Moniker, valName) {
				continue
			}
		}

		if isDelegateFlag == "0" {
			valInfo := DposValidatorInfo{
				ValidatorAddr:           validator.OperatorAddress,
				ValidatorName:           validator.Description.Moniker,
				ValidatorDelegateAmount: validator.Tokens,
				ValidatorOnlineRatio:    "",
				ValidatorCommision:      validator.Commission.CommissionRates.Rate,
				PersonDelegateAmount:    delVal,
			}

			valInfos = append(valInfos, valInfo)
			continue
		} else if isDelegateFlag == "1" {
			isDelegate = true
		} else if isDelegateFlag == "2" {
			isDelegate = false
		} else {
			return
		}

		if isDelegate == ok {
			valInfo := DposValidatorInfo{
				ValidatorAddr:           validator.OperatorAddress,
				ValidatorName:           validator.Description.Moniker,
				ValidatorDelegateAmount: validator.Tokens,
				ValidatorOnlineRatio:    "",
				ValidatorCommision:      validator.Commission.CommissionRates.Rate,
				PersonDelegateAmount:    delVal,
			}

			valInfos = append(valInfos, valInfo)
		}
	}

	
	if sortType == "1" {
		sort.Slice(valInfos, func(i, j int) bool {
			return valInfos[i].PersonDelegateAmount.LT(valInfos[j].PersonDelegateAmount)
		})
	}

	if sortType == "2" {
		sort.Slice(valInfos, func(i, j int) bool {
			return valInfos[i].PersonDelegateAmount.GT(valInfos[j].PersonDelegateAmount)
		})
	}

	if sortType == "3" {
		sort.Slice(valInfos, func(i, j int) bool {
			return valInfos[i].ValidatorCommision.LT(valInfos[j].ValidatorCommision)
		})
	}

	if sortType == "4" {
		sort.Slice(valInfos, func(i, j int) bool {
			return valInfos[i].ValidatorCommision.GT(valInfos[j].ValidatorCommision)
		})
	}

	
	start, end := client.Paginate(len(valInfos), int(page), int(limit), len(valInfos))
	if start < 0 || end < 0 {
		valInfos = []DposValidatorInfo{}
	} else {
		valInfos = valInfos[start:end]
	}

	resp.ValidatorList = valInfos
	resp.ValidatorAmount = int64(len(valInfos))

	return
}

type ValInfo struct {
	ValAddr           string `json:"val_addr"`            
	Name              string `json:"name"`                
	ValDelegateAmount string `json:"val_delegate_amount"` 
	SelfDelegateRate  string `json:"self_delegate_rate"`  
	MinSelfDelegate   string `json:"min_self_delegate"`   
	Identify          string `json:"identify"`            
	Status            int64  `json:"status"`              
	UnbondHeight      int64  `json:"unbond_height"`       
	UnbondTime        int64  `json:"unbond_time"`         
	Jail              bool   `json:"jail"`                
	Contact           string `json:"contact"`             
	DelegateAmount    string `json:"delegated"`           

	UpdateTime        int64  `json:"update_time"`         
	CommissionRate    string `json:"commission_rate"`     
	MaxCommissionRate string `json:"max_commission_rate"` 
	MaxChangeRate     string `json:"max_change_rate"`     

	BalanceAmount string `json:"balance_amount"`
	BalanceDenom  string `json:"balance_denom"`
}

func (this *DposClient) ValidatorDetail(from, valAddr string) (ValInfo, error) {

	var resp ValInfo
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	accFrom, err := sdk.AccAddressFromBech32(from)
	if err != nil {
		return resp, err
	}

	val, err := sdk.ValAddressFromBech32(valAddr)
	if err != nil {
		return resp, err
	}

	valAccAddress := sdk.AccAddress(val)

	queryValidatorInfoParams := stakingTypes.QueryValidatorParams{
		ValidatorAddr: val,
	}

	queryValidatorInfoParamsData, err := clientCtx.LegacyAmino.MarshalJSON(queryValidatorInfoParams)
	if err != nil {
		return resp, err
	}

	
	resBytes, _, err := clientCtx.QueryWithData("custom/staking/"+stakingTypes.QueryValidator, queryValidatorInfoParamsData)
	if err != nil {
		return resp, err
	}
	valInfo := stakingTypes.Validator{}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &valInfo)
	if err != nil {
		return resp, err
	}

	resp.ValAddr = valAddr
	resp.Name = valInfo.Description.Moniker
	resp.ValDelegateAmount = valInfo.BondedTokens().String()
	resp.MinSelfDelegate = valInfo.MinSelfDelegation.String()
	resp.Identify = valInfo.Description.Identity
	resp.Status = int64(valInfo.Status)
	resp.UnbondHeight = valInfo.UnbondingHeight
	resp.UnbondTime = valInfo.UnbondingTime.Unix()
	resp.Jail = valInfo.Jailed
	resp.Contact = valInfo.Description.SecurityContact
	resp.UpdateTime = valInfo.Commission.UpdateTime.Unix()
	resp.CommissionRate = core.RemoveStringLastZero(valInfo.Commission.Rate.String())
	resp.MaxCommissionRate = core.RemoveStringLastZero(valInfo.Commission.MaxChangeRate.String())
	resp.MaxChangeRate = core.RemoveStringLastZero(valInfo.Commission.MaxChangeRate.String())

	
	queryValSelfDelegateParams := stakingTypes.QueryDelegatorValidatorRequest{
		DelegatorAddr: valAccAddress.String(),
		ValidatorAddr: valAddr,
	}

	queryValSelfDelegateParamsData, err := clientCtx.LegacyAmino.MarshalJSON(queryValSelfDelegateParams)
	if err != nil {
		return resp, err
	}

	resBytes, _, err = clientCtx.QueryWithData("custom/staking/"+stakingTypes.QueryDelegation, queryValSelfDelegateParamsData)
	if err != nil {
		return resp, err
	}
	delegation := stakingTypes.DelegationResponse{}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &delegation)
	if err != nil {
		return resp, err
	}

	
	resp.SelfDelegateRate = core.RemoveStringLastZero(sdk.NewDecFromInt(delegation.Balance.Amount).Quo(sdk.NewDecFromInt(valInfo.BondedTokens())).String())

	
	queryFromSelfDelegateParams := stakingTypes.QueryDelegatorValidatorRequest{
		DelegatorAddr: accFrom.String(),
		ValidatorAddr: valAddr,
	}

	queryFromSelfDelegateParamsData, err := clientCtx.LegacyAmino.MarshalJSON(queryFromSelfDelegateParams)
	if err != nil {
		return resp, err
	}

	resBytes, _, err = util.QueryWithDataWithUnwrapErr(clientCtx, "custom/staking/"+stakingTypes.QueryDelegation, queryFromSelfDelegateParamsData)

	if err != nil && err.Error() != stakingTypes.ErrNoDelegation.Error() {
		return resp, err
	} else if err != nil && err.Error() == stakingTypes.ErrNoDelegation.Error() {
		resp.DelegateAmount = "0"
	} else {
		err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &delegation)
		if err != nil {
			return resp, err
		}

		resp.DelegateAmount = delegation.Balance.Amount.String()
	}

	
	
	queryBalanceParams := bankTypes.QueryBalanceRequest{
		Address: from,
		Denom:   core.GovDenom,
	}

	queryBalanceParamsData, err := clientCtx.LegacyAmino.MarshalJSON(queryBalanceParams)
	if err != nil {
		return resp, err
	}

	resBytes, _, err = util.QueryWithDataWithUnwrapErr(clientCtx, "custom/bank/"+bankTypes.QueryBalance, queryBalanceParamsData)

	balance := sdk.Coin{}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &balance)
	if err != nil {
		return resp, err
	}

	resp.BalanceAmount = balance.Amount.String()
	resp.BalanceDenom = balance.Denom

	return resp, nil
}


func (this *DposClient) QueryDelegatorDelegations(from string) (stakingTypes.DelegationResponses, error) {

	res := make([]stakingTypes.DelegationResponse, 0)
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	accFromAddr, err := sdk.AccAddressFromBech32(from)
	if err != nil {
		return res, err
	}

	requestDelegationsParams := stakingTypes.QueryDelegatorParams{
		DelegatorAddr: accFromAddr,
	}
	requestDelegationsParamsData, err := clientCtx.LegacyAmino.MarshalJSON(requestDelegationsParams)
	if err != nil {
		return res, err
	}

	resBytes, _, err := clientCtx.QueryWithData("custom/staking/"+stakingTypes.QueryDelegatorDelegations, requestDelegationsParamsData)
	if err != nil {
		return res, err
	}
	delegations := stakingTypes.DelegationResponses{}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &delegations)
	if err != nil {
		return res, err
	}

	return delegations, nil
}


func (this *DposClient) QueryDelegatorRewards(from string) (distributionTypes.QueryDelegatorTotalRewardsResponse, error) {

	res := distributionTypes.QueryDelegatorTotalRewardsResponse{}
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	accFromAddr, err := sdk.AccAddressFromBech32(from)
	if err != nil {
		return res, err
	}

	requestDelegationRewardsParams := distributionTypes.QueryDelegatorParams{
		DelegatorAddress: accFromAddr,
	}
	requestDelegationsParamsData, err := clientCtx.LegacyAmino.MarshalJSON(requestDelegationRewardsParams)
	if err != nil {
		return res, err
	}

	resBytes, _, err := clientCtx.QueryWithData("custom/distribution/"+distributionTypes.QueryDelegatorTotalRewards, requestDelegationsParamsData)
	if err != nil {
		return res, err
	}
	reward := distributionTypes.QueryDelegatorTotalRewardsResponse{}
	err = json.Unmarshal(resBytes, &reward)
	if err != nil {
		return res, err
	}

	return reward, nil
}

func (this *DposClient) QueryDelegatorValidators(from string) (sdk.Coin, uint64, sdk.DecCoin, error) {
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	accFromAddr, err := sdk.AccAddressFromBech32(from)
	if err != nil {
		return sdk.Coin{}, 0, sdk.DecCoin{}, err
	}

	requestDelegationValidatorsParams := stakingTypes.QueryDelegatorParams{
		DelegatorAddr: accFromAddr,
	}
	requestDelegationValidatorsParamsData, err := clientCtx.LegacyAmino.MarshalJSON(requestDelegationValidatorsParams)
	if err != nil {
		return sdk.Coin{}, 0, sdk.DecCoin{}, err
	}

	resBytes, _, err := clientCtx.QueryWithData("custom/staking/"+stakingTypes.QueryDelegatorValidators, requestDelegationValidatorsParamsData)
	if err != nil {
		return sdk.Coin{}, 0, sdk.DecCoin{}, err
	}
	validators := stakingTypes.Validators{}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &validators)
	if err != nil {
		return sdk.Coin{}, 0, sdk.DecCoin{}, err
	}

	msgs := make([]sdk.Msg, 0)
	for _, validator := range validators {
		valAddress, err := sdk.ValAddressFromBech32(validator.OperatorAddress)
		if err != nil {
			return sdk.Coin{}, 0, sdk.DecCoin{}, err
		}

		params := distributionTypes.NewQueryDelegationRewardsParams(accFromAddr, valAddress)
		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		rewardsBytes, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", distributionTypes.QuerierRoute, distributionTypes.QueryDelegationRewards), bz)
		if err != nil {
			return sdk.Coin{}, 0, sdk.DecCoin{}, err
		}

		rewards := sdk.DecCoins{}
		err = clientCtx.LegacyAmino.UnmarshalJSON(rewardsBytes, &rewards)
		if err != nil {
			return sdk.Coin{}, 0, sdk.DecCoin{}, err
		}

		
		if len(rewards) > 0 {
			rewardMsg := distributionTypes.NewMsgWithdrawDelegatorReward(accFromAddr, valAddress)
			msgs = append(msgs, rewardMsg)
		}
	}

	
	commisionMsgs, err := this.GetWithdrawDelegatorRewardsMsg(from)
	if err != nil {
		return sdk.Coin{}, 0, sdk.DecCoin{}, err
	}

	if commisionMsgs != nil {
		msgs = append(msgs, commisionMsgs)
	}

	seq, err := this.TxClient.FindAccountNumberSeq(from)
	if err != nil {
		return sdk.Coin{}, 0, sdk.DecCoin{}, err
	}
	fee, gas, gasPrice, _, err := this.TxClient.GasInfo(seq, msgs...)
	if err != nil {
		return sdk.Coin{}, 0, sdk.DecCoin{}, err
	}
	return fee, gas, gasPrice, nil
}

func (this *DposClient) GetWithdrawDelegatorRewardsMsg(addr string) (sdk.Msg, error) {
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	accFromAddress, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return nil, err
	}

	accValaddress := sdk.ValAddress(accFromAddress)
	if err != nil {
		return nil, err
	}

	
	commisionParams := distributionTypes.NewQueryValidatorCommissionParams(accValaddress)
	bz, err := clientCtx.LegacyAmino.MarshalJSON(commisionParams)
	if err != nil {
		return nil, err
	}

	commisionGrpcRes, _, err := clientCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", distributionTypes.QuerierRoute, distributionTypes.QueryValidatorCommission),
		bz,
	)
	if err != nil {
		return nil, err
	}

	var commision distributionTypes.ValidatorAccumulatedCommission
	err = clientCtx.LegacyAmino.UnmarshalJSON(commisionGrpcRes, &commision)
	if err != nil {
		return nil, err
	}

	
	if len(commision.Commission) > 0 {
		commsionMsg := distributionTypes.NewMsgWithdrawValidatorCommission(accValaddress)
		return commsionMsg, nil
	}

	return nil, nil
}
