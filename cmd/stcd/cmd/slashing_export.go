package cmd

import (
	"encoding/json"
	slashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"time"
)

func ExportSlashings(jsonData json.RawMessage) ([]byte, error) {
	slashingState := slashingTypes.GenesisState{}
	encCfg.Codec.MustUnmarshalJSON(jsonData, &slashingState)

	
	slashingState.MissedBlocks = make([]slashingTypes.ValidatorMissedBlocks, 0)

	
	newSigningInfos := make([]slashingTypes.SigningInfo, 0)
	for _, signingInfo := range slashingState.SigningInfos {
		if signingInfo.Address == "dstvalcons13h84n3qt5653gflu4syzg3laz43dl9emnm8nyy" {
			continue
		}
		newSigningInfos = append(newSigningInfos, slashingTypes.SigningInfo{
			Address: signingInfo.Address,
			ValidatorSigningInfo: slashingTypes.ValidatorSigningInfo{
				Address:             signingInfo.ValidatorSigningInfo.Address,
				StartHeight:         0,
				IndexOffset:         0,
				JailedUntil:         time.Unix(0, 0),
				Tombstoned:          false,
				MissedBlocksCounter: 0,
			},
		})

	}

	slashingState.SigningInfos = newSigningInfos

	
	slashingState.Params.SignedBlocksWindow = 100

	slashingByte := encCfg.Codec.MustMarshalJSON(&slashingState)
	return slashingByte, nil
}
