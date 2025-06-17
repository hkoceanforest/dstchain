package cmd

import (
	"encoding/json"
	daoTypes "freemasonry.cc/blockchain/x/dao/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func ExportDao(jsonData json.RawMessage) ([]byte, error) {
	daoState := daoTypes.GenesisState{}
	encCfg.Codec.MustUnmarshalJSON(jsonData, &daoState)

	
	if daoState.Params.ClusterLevels[0].Level == 1 {
		addClusterLevels := []daoTypes.ClusterLevel{
			{
				Level:        0,
				BurnAmount:   sdk.ZeroInt(),
				MemberAmount: 1,
			},
		}

		daoState.Params.ClusterLevels = append(addClusterLevels, daoState.Params.ClusterLevels...)
	}

	daoState.Params.CrossFee = daoTypes.DefaulIBCTransferFee()
	daoState.Params.MintBlockInterval = 3600

	
	amount := sdk.MustNewDecFromStr("10000000000000000000").TruncateInt()
	daoState.Params.ClusterLevels[1].BurnAmount = amount
	daoState.Params.ClusterLevels[1].DaoLimit = amount
	
	for n, _ := range daoState.Clusters {
		daoState.Clusters[n].ClusterPowerMembers = make(map[string]daoTypes.ClusterPowerMemberExport)
		daoState.Clusters[n].ClusterPower = sdk.ZeroDec()
		daoState.Clusters[n].ClusterBurnAmount = sdk.ZeroDec()
	}

	
	for m, _ := range daoState.PersonalClusters {
		daoState.PersonalClusters[m].FreezePower = sdk.ZeroDec()

		daoState.PersonalClusters[m].ActivePower = sdk.ZeroDec()
		daoState.PersonalClusters[m].AllBurn = sdk.ZeroDec()
	}

	
	newDaoClusterLevels := daoState.Params.ClusterLevels
	newDaoClusterLevels[0] = daoTypes.ClusterLevel{
		Level:        1,
		BurnAmount:   sdk.MustNewDecFromStr("10000000000000000").TruncateInt(),
		MemberAmount: 1,
	}

	newDaoClusterLevels = append([]daoTypes.ClusterLevel{
		{
			Level:        0,
			BurnAmount:   sdk.NewInt(0),
			MemberAmount: 1,
		},
	}, newDaoClusterLevels...)

	daoState.Params.ClusterLevels = newDaoClusterLevels

	daoByte := encCfg.Codec.MustMarshalJSON(&daoState)
	return daoByte, nil
}
