package cmd

import (
	"encoding/json"
	"freemasonry.cc/blockchain/core"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"time"
)

func ExportStaking(jsonData json.RawMessage) ([]byte, sdk.Int, map[string]sdk.Int, []stakingTypes.Validator, error) {
	stakingState := stakingTypes.GenesisState{}
	encCfg.Codec.MustUnmarshalJSON(jsonData, &stakingState)

	
	totalStaked := sdk.ZeroInt()

	
	newdelegations := make([]stakingTypes.Delegation, 0)

	minStakeMap := make(map[string]sdk.Int)

	genesisStakeAmount, ok := sdk.NewIntFromString("10000000000000000000000")
	if !ok {
		panic("ExportStaking genesisStakeAmount error")
	}

	normalStakeAmount, ok := sdk.NewIntFromString("1000000000000000000")
	if !ok {
		panic("ExportStaking normalStakeAmount error")
	}

	for k, validator := range stakingState.Validators {
		if validator.OperatorAddress == "dstvaloper1fcaqgl2fxg65tspkt27muwvu2t3uafelhnjahp" {
			continue
		}

		if validator.OperatorAddress == "dstvaloper1ll30h0xykgduvxxfnpy4h6yzl0770pgn7hn3lz" {
			stakingState.Validators[k].MinSelfDelegation = genesisStakeAmount
			validator.MinSelfDelegation = genesisStakeAmount
			validator.Status = 3
		} else {
			stakingState.Validators[k].MinSelfDelegation = normalStakeAmount
			validator.MinSelfDelegation = normalStakeAmount
			validator.Status = 1
		}

		
		if newValOper, ok := core.ValReplace[validator.OperatorAddress]; ok {
			
			valAddr, err := sdk.ValAddressFromBech32(newValOper)
			if err != nil {
				return nil, sdk.ZeroInt(), nil, nil, err
			}
			accAddr := sdk.AccAddress(valAddr.Bytes())

			
			newdelegations = append(newdelegations, stakingTypes.Delegation{
				DelegatorAddress: accAddr.String(),
				ValidatorAddress: newValOper,
				Shares:           sdk.NewDecFromInt(validator.MinSelfDelegation),
			})
		} else {
			
			valAddr, err := sdk.ValAddressFromBech32(validator.OperatorAddress)
			if err != nil {
				return nil, sdk.ZeroInt(), nil, nil, err
			}
			accAddr := sdk.AccAddress(valAddr.Bytes())

			newdelegations = append(newdelegations, stakingTypes.Delegation{
				DelegatorAddress: accAddr.String(),
				ValidatorAddress: validator.OperatorAddress,
				Shares:           sdk.NewDecFromInt(validator.MinSelfDelegation),
			})
		}
		if validator.Status == 3 {
			totalStaked = totalStaked.Add(validator.MinSelfDelegation)
		}

		minStakeMap[validator.OperatorAddress] = validator.MinSelfDelegation

	}
	stakingState.Delegations = newdelegations

	
	newValidators := make([]stakingTypes.Validator, 0)
	for _, validator := range stakingState.Validators {
		
		if validator.OperatorAddress == "dstvaloper1fcaqgl2fxg65tspkt27muwvu2t3uafelhnjahp" {
			continue
		}

		if validator.OperatorAddress == "dstvaloper1ll30h0xykgduvxxfnpy4h6yzl0770pgn7hn3lz" {
			newValidators = append(newValidators, stakingTypes.Validator{
				OperatorAddress:   validator.OperatorAddress,
				ConsensusPubkey:   validator.ConsensusPubkey,
				Jailed:            false,
				Status:            3,
				Tokens:            validator.MinSelfDelegation,
				DelegatorShares:   sdk.NewDecFromInt(validator.MinSelfDelegation),
				Description:       validator.Description,
				UnbondingHeight:   0,
				UnbondingTime:     time.Time{},
				Commission:        validator.Commission,
				MinSelfDelegation: validator.MinSelfDelegation,
			})
		} else {

			newValidators = append(newValidators, stakingTypes.Validator{
				OperatorAddress:   validator.OperatorAddress,
				ConsensusPubkey:   validator.ConsensusPubkey,
				Jailed:            false,
				Status:            1,
				Tokens:            validator.MinSelfDelegation,
				DelegatorShares:   sdk.NewDecFromInt(validator.MinSelfDelegation),
				Description:       validator.Description,
				UnbondingHeight:   0,
				UnbondingTime:     time.Time{},
				Commission:        validator.Commission,
				MinSelfDelegation: validator.MinSelfDelegation,
			})
		}

	}
	stakingState.Validators = newValidators

	
	stakingState.Redelegations = make([]stakingTypes.Redelegation, 0)
	
	stakingState.UnbondingDelegations = make([]stakingTypes.UnbondingDelegation, 0)
	
	stakingState.LastTotalPower = sdk.NewInt(10000)
	stakingState.LastValidatorPowers = make([]stakingTypes.LastValidatorPower, 0)
	stakingState.LastValidatorPowers = append(stakingState.LastValidatorPowers, stakingTypes.LastValidatorPower{
		Address: "dstvaloper1ll30h0xykgduvxxfnpy4h6yzl0770pgn7hn3lz",
		Power:   10000,
	})
	
	
	
	
	
	
	
	

	distByte := encCfg.Codec.MustMarshalJSON(&stakingState)
	return distByte, totalStaked, minStakeMap, stakingState.Validators, nil
}
