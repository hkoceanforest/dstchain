package contracts

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

var (
	//go:embed compiled_contracts/GenesisIdo.json
	genesisIdoJSON []byte

	GenesisIdoNContract evmtypes.CompiledContract
)

func init() {
	err := json.Unmarshal(genesisIdoJSON, &GenesisIdoNContract)
	if err != nil {
		panic(err)
	}
}
