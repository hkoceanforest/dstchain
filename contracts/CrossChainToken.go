package contracts

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

var (
	//go:embed compiled_contracts/CrossChainToken.json
	crossChainTokenJSON []byte

	CrossChainTokenJSONContract evmtypes.CompiledContract
)

func init() {
	err := json.Unmarshal(crossChainTokenJSON, &CrossChainTokenJSONContract)
	if err != nil {
		panic(err)
	}
}
