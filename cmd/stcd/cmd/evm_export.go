package cmd

import (
	"encoding/json"
	"github.com/evmos/ethermint/x/evm/types"
)

func ExportEvm(jsonData json.RawMessage) ([]byte, error) {
	evmState := types.GenesisState{}
	encCfg.Codec.MustUnmarshalJSON(jsonData, &evmState)

	evmState.Accounts = make([]types.GenesisAccount, 0)

	govByte := encCfg.Codec.MustMarshalJSON(&evmState)
	return govByte, nil
}
