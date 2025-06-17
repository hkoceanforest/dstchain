package contracts

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

var (
	//go:embed compiled_contracts/LP.json
	LPJSON []byte

	LPContract evmtypes.CompiledContract
)

func init() {
	err := json.Unmarshal(LPJSON, &LPContract)
	if err != nil {
		panic(err)
	}
}
