package cmd

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)


func ExportDistribution(jsonData json.RawMessage, stakeMap map[string]sdk.Int) ([]byte, error) {
	distributionState := distTypes.GenesisState{}
	encCfg.Codec.MustUnmarshalJSON(jsonData, &distributionState)

	
	delegatorStartingInfoRecord := make([]distTypes.DelegatorStartingInfoRecord, 0)
	for _, delegatorStartingInfo := range distributionState.DelegatorStartingInfos {
		
		if delegatorStartingInfo.ValidatorAddress == "dstvaloper1fcaqgl2fxg65tspkt27muwvu2t3uafelhnjahp" {
			continue
		}

		if delegatorStartingInfo.DelegatorAddress == "dst1ll30h0xykgduvxxfnpy4h6yzl0770pgneerr6w" {
			delegatorStartingInfo.DelegatorAddress = "dst1kxcnfvtp6ep042vudctgda29cywkpzp0ttky75"
		}

		if delegatorStartingInfo.StartingInfo.PreviousPeriod == 1 {
			delegatorStartingInfoRecord = append(delegatorStartingInfoRecord, distTypes.DelegatorStartingInfoRecord{
				DelegatorAddress: delegatorStartingInfo.DelegatorAddress,
				ValidatorAddress: delegatorStartingInfo.ValidatorAddress,
				StartingInfo: distTypes.DelegatorStartingInfo{
					PreviousPeriod: 1,
					Stake:          sdk.NewDecFromInt(stakeMap[delegatorStartingInfo.ValidatorAddress]),
					Height:         0,
				},
			})
		}
	}
	distributionState.DelegatorStartingInfos = delegatorStartingInfoRecord

	
	newDelegatorWithdrawInfos := make([]distTypes.DelegatorWithdrawInfo, 0)
	distributionState.DelegatorWithdrawInfos = newDelegatorWithdrawInfos

	
	newFeePool := distTypes.InitialFeePool()
	distributionState.FeePool = newFeePool

	
	newOutstandingRewards := make([]distTypes.ValidatorOutstandingRewardsRecord, 0)
	distributionState.OutstandingRewards = newOutstandingRewards

	
	newValidatorAccumulatedCommissions := make([]distTypes.ValidatorAccumulatedCommissionRecord, 0)
	for _, validatorAccumulatedCommission := range distributionState.ValidatorAccumulatedCommissions {
		
		if validatorAccumulatedCommission.ValidatorAddress == "dstvaloper1fcaqgl2fxg65tspkt27muwvu2t3uafelhnjahp" {
			continue
		}

		newValidatorAccumulatedCommissions = append(newValidatorAccumulatedCommissions, distTypes.ValidatorAccumulatedCommissionRecord{
			ValidatorAddress: validatorAccumulatedCommission.ValidatorAddress,
			Accumulated: distTypes.ValidatorAccumulatedCommission{
				Commission: sdk.NewDecCoins(sdk.NewDecCoin("nxn", sdk.ZeroInt()), sdk.NewDecCoin("dst", sdk.ZeroInt())),
			},
		})
	}
	distributionState.ValidatorAccumulatedCommissions = newValidatorAccumulatedCommissions

	
	newValidatorCurrentRewards := make([]distTypes.ValidatorCurrentRewardsRecord, 0)
	for _, validatorCurrentRewards := range distributionState.ValidatorCurrentRewards {
		
		if validatorCurrentRewards.ValidatorAddress == "dstvaloper1fcaqgl2fxg65tspkt27muwvu2t3uafelhnjahp" {
			continue
		}

		newValidatorCurrentRewards = append(newValidatorCurrentRewards, distTypes.ValidatorCurrentRewardsRecord{
			ValidatorAddress: validatorCurrentRewards.ValidatorAddress,
			Rewards: distTypes.ValidatorCurrentRewards{
				Rewards: sdk.DecCoins{},
				Period:  2,
			},
		})

	}

	distributionState.ValidatorCurrentRewards = newValidatorCurrentRewards

	

	historyMap := make(map[string]struct{})

	newValidatorHistoricalRewards := make([]distTypes.ValidatorHistoricalRewardsRecord, 0)
	for _, validatorHistoricalRewards := range distributionState.ValidatorHistoricalRewards {
		
		if validatorHistoricalRewards.ValidatorAddress == "dstvaloper1fcaqgl2fxg65tspkt27muwvu2t3uafelhnjahp" {
			continue
		}

		if _, ok := historyMap[validatorHistoricalRewards.ValidatorAddress]; ok {
			continue
		}

		newValidatorHistoricalRewards = append(newValidatorHistoricalRewards, distTypes.ValidatorHistoricalRewardsRecord{
			ValidatorAddress: validatorHistoricalRewards.ValidatorAddress,
			Period:           1,
			Rewards: distTypes.ValidatorHistoricalRewards{
				CumulativeRewardRatio: sdk.DecCoins{},
				ReferenceCount:        2,
			},
		})

		historyMap[validatorHistoricalRewards.ValidatorAddress] = struct{}{}

	}
	distributionState.ValidatorHistoricalRewards = newValidatorHistoricalRewards

	
	newValidatorSlashEvents := make([]distTypes.ValidatorSlashEventRecord, 0)
	distributionState.ValidatorSlashEvents = newValidatorSlashEvents

	
	distributionStateExport, err := encCfg.Codec.MarshalJSON(&distributionState)
	if err != nil {
		return nil, err
	}
	return distributionStateExport, nil
}
