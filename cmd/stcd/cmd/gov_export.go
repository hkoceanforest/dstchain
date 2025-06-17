package cmd

import (
	"encoding/json"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"time"
)

func ExportGov(jsonData json.RawMessage) ([]byte, error) {
	govState := govTypes.GenesisState{}
	encCfg.Codec.MustUnmarshalJSON(jsonData, &govState)

	govState.Proposals = nil

	govState.StartingProposalId = 1

	votingPeriod := time.Duration(4) * time.Hour

	govState.VotingParams.VotingPeriod = &votingPeriod

	govByte := encCfg.Codec.MustMarshalJSON(&govState)
	return govByte, nil
}
