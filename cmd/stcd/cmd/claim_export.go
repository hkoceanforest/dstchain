package cmd

import (
	"encoding/json"
	"github.com/evmos/evmos/v10/x/claims/types"
)

func ExportClaim(jsonData json.RawMessage) ([]byte, error) {
	claimState := types.GenesisState{}
	encCfg.Codec.MustUnmarshalJSON(jsonData, &claimState)

	claimState.Params.EVMChannels = []string{
		"channel-0",
	}

	claimByte := encCfg.Codec.MustMarshalJSON(&claimState)
	return claimByte, nil
}
