package contracts

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

var (
	//go:embed compiled_contracts/Usdt.json
	usdtJSON []byte

	UsdtContract evmtypes.CompiledContract
)

func init() {
	err := json.Unmarshal(usdtJSON, &UsdtContract)
	if err != nil {
		panic(err)
	}
}
