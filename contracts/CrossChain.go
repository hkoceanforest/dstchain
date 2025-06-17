package contracts

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

var (
	//go:embed compiled_contracts/CrossChain.json
	crossChainJSON []byte

	CrossChainJSONContract evmtypes.CompiledContract
)

func init() {
	err := json.Unmarshal(crossChainJSON, &CrossChainJSONContract)
	if err != nil {
		panic(err)
	}
}
