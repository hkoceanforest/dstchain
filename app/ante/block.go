package ante

import sdk "github.com/cosmos/cosmos-sdk/types"

func BlockGasLimit(ctx sdk.Context) uint64 {
	blockGasMeter := ctx.BlockGasMeter()

	
	if blockGasMeter != nil && blockGasMeter.Limit() != 0 {
		return blockGasMeter.Limit()
	}

	
	cp := ctx.ConsensusParams()
	if cp == nil || cp.Block == nil {
		return 0
	}

	maxGas := cp.Block.MaxGas
	return uint64(maxGas)
}
