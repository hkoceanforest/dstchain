package ante

import (
	"freemasonry.cc/blockchain/core"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibcTypes "github.com/cosmos/ibc-go/v5/modules/apps/transfer/types"
)

type IBCFeeDecorator struct {
	daoKeeper      DaoKeeper
	contractKeeper ContractsKeeper
}

func NewIBCFeeDecorator(ck ContractsKeeper, dk DaoKeeper) IBCFeeDecorator {
	return IBCFeeDecorator{contractKeeper: ck, daoKeeper: dk}
}

func (ibcFeeDecorator IBCFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	for _, msg := range tx.GetMsgs() {
		switch m := msg.(type) {
		case *ibcTypes.MsgTransfer:
			
			if m.Token.Denom == core.GovDenom {
				
				return next(ctx, tx, simulate)
			}
			pairID := ibcFeeDecorator.contractKeeper.GetTokenPairID(ctx, m.Token.Denom)
			if len(pairID) == 0 {
				
				
				return next(ctx, tx, simulate)
			}

			pair, _ := ibcFeeDecorator.contractKeeper.GetTokenPair(ctx, pairID)
			if !pair.Enabled {
				
				return next(ctx, tx, simulate)
			}
			err := ibcFeeDecorator.daoKeeper.DeductionToken(ctx, pair.Erc20Address, m.Sender)
			if err != nil {
				return ctx, err
			}
		}

	}
	return next(ctx, tx, simulate)
}
