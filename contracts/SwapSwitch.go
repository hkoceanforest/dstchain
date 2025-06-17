package contracts

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

var (
	//go:embed compiled_contracts/SwapSwitch.json
	SwapSwitchJSON []byte

	SwapSwitchJSONContract evmtypes.CompiledContract
)

func init() {
	err := json.Unmarshal(SwapSwitchJSON, &SwapSwitchJSONContract)
	if err != nil {
		panic(err)
	}
}
