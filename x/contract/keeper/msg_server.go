package keeper

import (
	"context"
	"encoding/json"
	"errors"
	"freemasonry.cc/blockchain/contracts"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/x/contract/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ethereum/go-ethereum/common"
	types2 "github.com/evmos/ethermint/x/evm/types"
	erc20types "github.com/evmos/evmos/v10/x/erc20/types"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmtypes "github.com/tendermint/tendermint/types"
	"math/big"
	"strconv"
)

var _ types.MsgServer = &msgServer{}

type msgServer struct {
	Keeper
	logPrefix string
}

func (s msgServer) RegisterErc20(goCtx context.Context, msg *types.MsgRegisterErc20) (*types.MsgEmptyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	signers := msg.GetSigners()
	if !signers[0].Equals(s.accountKeeper.GetModuleAddress(govtypes.ModuleName)) {
		return &types.MsgEmptyResponse{}, sdkerrors.Wrapf(govtypes.ErrInvalidSigner, signers[0].String())
	}
	pair, err := s.RegisterERC20(ctx, msg.ContractAddress, msg.Denom, msg.Owner)
	if err != nil {
		return &types.MsgEmptyResponse{}, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			erc20types.EventTypeRegisterERC20,
			sdk.NewAttribute(erc20types.AttributeKeyCosmosCoin, pair.Denom),
			sdk.NewAttribute(erc20types.AttributeKeyERC20Token, pair.Erc20Address),
		),
	)
	return &types.MsgEmptyResponse{}, nil
}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper, logPrefix: "contract | msgServer | "}
}

func (s msgServer) AppTokenIssue(goCtx context.Context, msg *types.MsgAppTokenIssue) (*types.MsgEmptyResponse, error) {

	log := core.BuildLog(core.GetStructFuncName(s), core.LmChainMsgServer)

	log.Info("tokenDecimals:", msg.Decimals)

	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, _ := sdk.AccAddressFromBech32(msg.FromAddress)
	from := common.BytesToAddress(addr)

	contractAddr := common.HexToAddress(s.GetTokenFactoryContractAddress(ctx))

	preMintAmount, err := String2BigInt(msg.PreMintAmount)
	if err != nil {
		return nil, err
	}

	tokenDecimalsInt64, err := strconv.ParseInt(msg.Decimals, 10, 8)
	if err != nil {
		return nil, err
	}

	for i := int64(1); i <= tokenDecimalsInt64; i++ {
		preMintAmount.Mul(preMintAmount, big.NewInt(10))
	}

	param := types.TokenFactoryCreateParams{
		from,
		msg.Name,
		msg.Symbol,
		preMintAmount,
		uint8(tokenDecimalsInt64),
	}

	stakeFactory := contracts.AppTokenIssueJSONContract
	resp1, err := s.CallEVM(ctx, stakeFactory.ABI, from, contractAddr, true, "createToken", param)
	if err != nil {
		log.Error("AppTokenIssue CallEVM Error:", err.Error())
		return nil, err
	}

	if resp1.VmError != "" {
		log.Error("VmError Error:", resp1.VmError)
		return nil, errors.New(resp1.VmError)
	}

	attrs := []sdk.Attribute{
		sdk.NewAttribute(sdk.AttributeKeyAmount, ""),
		
		sdk.NewAttribute(types2.AttributeKeyEthereumTxHash, resp1.Hash),
		
		sdk.NewAttribute(types2.AttributeKeyTxIndex, strconv.FormatUint(0, 10)),
		
		sdk.NewAttribute(types2.AttributeKeyTxGasUsed, strconv.FormatUint(resp1.GasUsed, 10)),
	}

	if len(ctx.TxBytes()) > 0 {
		
		hash := tmbytes.HexBytes(tmtypes.Tx(ctx.TxBytes()).Hash())
		attrs = append(attrs, sdk.NewAttribute(types2.AttributeKeyTxHash, hash.String()))
	}

	if resp1.Failed() {
		attrs = append(attrs, sdk.NewAttribute(types2.AttributeKeyEthereumTxFailed, resp1.VmError))
	}

	

	resp, err := s.CallEVM(
		ctx,
		stakeFactory.ABI,
		from,
		contractAddr,
		false,
		"getUserNewCreate",
		from,
	)
	if err != nil {
		log.Error("Query Token Issue Error:", err.Error())
		return nil, err
	}

	var interfaceArr []interface{}

	if interfaceArr, err = stakeFactory.ABI.Unpack("getUserNewCreate", resp.Ret); err != nil {
		log.Error("Query Token Issue Unpack Error:", err.Error())
		return nil, err
	}

	tokenInfoJson, _ := json.Marshal(interfaceArr[0])

	var tokenInfo types.TokenNewCreateInfo
	err = json.Unmarshal(tokenInfoJson, &tokenInfo)
	if err != nil {
		return nil, err
	}

	if tokenInfo.Symbol != msg.Symbol {
		return nil, errors.New("Error Token Issue")
	}

	
	txLogAttrs := make([]sdk.Attribute, len(resp1.Logs))
	for i, log := range resp1.Logs {
		value, err := json.Marshal(log)
		if err != nil {
			return nil, sdkerrors.Wrap(err, "failed to encode log")
		}
		txLogAttrs[i] = sdk.NewAttribute(types2.AttributeKeyTxLog, string(value))
	}

	
	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeAppTokenIssue,

				
				sdk.NewAttribute(types.EventTypeFromAddress, msg.FromAddress),
				
				sdk.NewAttribute(types.EventTypeETHFromAddress, from.String()),

				
				sdk.NewAttribute(types.EventTypeTokenAmount, msg.PreMintAmount),

				
				sdk.NewAttribute(types.EventTypeTokenSymbol, msg.Symbol),
				
				sdk.NewAttribute(types.EventTypeTokenName, msg.Name),
				
				sdk.NewAttribute(types.EventTypeTokenDecimal, msg.Decimals),
				
				sdk.NewAttribute(types.EventTypeTokenLogo, msg.LogoUrl),
				
				sdk.NewAttribute(types.EventTypeContractAddress, tokenInfo.TokenAddress.String()),
			),
			sdk.NewEvent(
				types2.EventTypeEthereumTx,
				attrs...,
			),
			sdk.NewEvent(
				types2.EventTypeTxLog,
				txLogAttrs...,
			),
		},
	)

	return nil, nil
}

func String2BigInt(s string) (*big.Int, error) {
	n := big.NewInt(0)
	_, ok := n.SetString(s, 10)
	if !ok {
		return nil, core.ErrStringNumber
	}
	return n, nil
}
