package contracts

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

var (
	//go:embed compiled_contracts/Multicall.json
	MulticallJSON []byte

	MulticallContract evmtypes.CompiledContract
)

func init() {
	err := json.Unmarshal(MulticallJSON, &MulticallContract)
	if err != nil {
		panic(err)
	}
}
