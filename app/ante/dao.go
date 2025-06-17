package ante

import (
	errorsmod "cosmossdk.io/errors"
	"freemasonry.cc/blockchain/contracts"
	"freemasonry.cc/blockchain/core"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authAnt "github.com/cosmos/cosmos-sdk/x/auth/ante"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"math/big"
	"strings"
)

type DaoAccountVerificationDecorator struct {
	ak        authAnt.AccountKeeper
	ck        ContractsKeeper
	daoKeeper DaoKeeper
}

func NewDaoAccountVerificationDecorator(ak authAnt.AccountKeeper, dk DaoKeeper, ck ContractsKeeper) DaoAccountVerificationDecorator {
	return DaoAccountVerificationDecorator{
		ak:        ak,
		daoKeeper: dk,
		ck:        ck,
	}
}

func (davd DaoAccountVerificationDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (newCtx sdk.Context, err error) {
	if !ctx.IsCheckTx() {
		return next(ctx, tx, simulate)
	}

	for _, msg := range tx.GetMsgs() {
		

		msgEthTx, ok := msg.(*evmtypes.MsgEthereumTx)
		if ok {
			msgEt := msgEthTx.AsTransaction()

			
			if msgEt.To() != nil && strings.ToLower(msgEt.To().String()) == strings.ToLower(davd.ck.GetRedPacketContractAddress(ctx)) {
				redPacketAbi := contracts.RedPacketContract.ABI

				orignalData := msgEthTx.AsTransaction().Data()

				if len(orignalData) >= 5 {
					methodData := orignalData[4:]

					m, err := redPacketAbi.MethodById(orignalData[:4])
					if err != nil {
						return ctx, err
					}

					if m.Name == "sendBag" {
						paramsList, err := redPacketAbi.Methods["sendBag"].Inputs.Unpack(methodData)
						if err != nil {
							return ctx, err
						}
						sender := msg.GetSigners()[0]
						clusterId := paramsList[0].(string)
						
						cluster, err := davd.daoKeeper.GetClusterByChatId(ctx, clusterId)
						if err != nil {
							
							return ctx, errorsmod.Wrap(
								err,
								core.GetClusterNotFound.Error(),
							)
						}

						if _, ok := cluster.ClusterDeviceMembers[sender.String()]; !ok {
							
							return ctx, core.ErrNotInCluster
						}

					}

					if m.Name == "openBag" {
						paramsList, err := redPacketAbi.Methods["openBag"].Inputs.Unpack(methodData)
						if err != nil {
							return ctx, err
						}

						sender := msg.GetSigners()[0]

						
						rid := paramsList[0].(*big.Int)

						
						clusterId, err := davd.daoKeeper.GetRidClusterId(ctx, rid.String())
						if err != nil {
							return ctx, err
						}

						
						cluster, err := davd.daoKeeper.GetCluster(ctx, clusterId)
						if err != nil {
							return ctx, err
						}

						
						if _, ok := cluster.ClusterDeviceMembers[sender.String()]; !ok {
							return ctx, core.ErrNotInCluster
						}

					}

				}

			}

		}
	}
	return next(ctx, tx, simulate)
}
