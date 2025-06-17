package contracts

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

var (
	//go:embed compiled_contracts/GenesisIdoNolimit.json
	genesisIdoNoLimitJSON []byte

	GenesisIdoNoLimitContract evmtypes.CompiledContract
)

func init() {
	err := json.Unmarshal(genesisIdoNoLimitJSON, &GenesisIdoNoLimitContract)
	if err != nil {
		panic(err)
	}
}
