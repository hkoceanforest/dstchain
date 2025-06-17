package contracts

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

var (
	//go:embed compiled_contracts/Wdst.json
	wdstJSON []byte

	WdstJSONContract evmtypes.CompiledContract
)

func init() {
	err := json.Unmarshal(wdstJSON, &WdstJSONContract)
	if err != nil {
		panic(err)
	}
}
