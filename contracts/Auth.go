package contracts

import (
	_ "embed"
	"encoding/json"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

var (
	//go:embed compiled_contracts/Auth.json
	authJSON []byte

	AuthJSONContract evmtypes.CompiledContract
)

func init() {
	err := json.Unmarshal(authJSON, &AuthJSONContract)
	if err != nil {
		panic(err)
	}
}
