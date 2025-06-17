package keeper

import (
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/x/contract/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	contracts2 "github.com/evmos/evmos/v10/contracts"
	erc20types "github.com/evmos/evmos/v10/x/erc20/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		var (
			res []byte
			err error
		)
		switch path[0] {
		case types.QueryParams:
			return queryParams(ctx, req, k, legacyQuerierCdc)
		case types.QueryContractCode:
			return queryContractCode(ctx, req, k, legacyQuerierCdc)
		case types.QueryTokenPair:
			return queryTokenPair(ctx, req, k, legacyQuerierCdc)
		case types.QueryMainTokenBalances:
			return queryMainTokenBalances(ctx, req, k, legacyQuerierCdc)
		default:
			err = sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint: %s", types.ModuleName, path[0])
		}

		return res, err
	}
}

type MainTokenBalance struct {
	Denom  string  `json:"denom"`
	Amount sdk.Int `json:"amount"`
}

func queryMainTokenBalances(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainKeeper)

	
	var accAccount sdk.AccAddress
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &accAccount)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	evmAccount := common.BytesToAddress(accAccount.Bytes())

	
	allBalances := k.BankKeeper.GetAllBalances(ctx, accAccount)

	
	newAllBalances := make([]sdk.Coin, 0)
	for _, coin := range allBalances {
		if coin.Denom != core.UsdtDenom {
			newAllBalances = append(newAllBalances, coin)
		}
	}

	allBalances = sdk.NewCoins(newAllBalances...)

	allMainTokens := map[string]bool{
		core.BaseDenom: false,
		core.GovDenom:  false,
	}

	for _, coin := range allBalances {
		if _, ok := allMainTokens[coin.Denom]; ok {
			allMainTokens[coin.Denom] = true
		}
	}

	for mainToken, isExist := range allMainTokens {
		if isExist == false {
			newToken := sdk.NewCoin(mainToken, sdk.ZeroInt())
			allBalances = append(allBalances, newToken)
		}
	}

	
	usdtBalanceInt := sdk.ZeroInt()
	queryUsdtParams := erc20types.QueryTokenPairRequest{
		Token: core.UsdtDenom,
	}
	ctxs := sdk.WrapSDKContext(ctx)
	resp, err := k.erc20Keeper.TokenPair(ctxs, &queryUsdtParams)
	if err != nil {
		usdtBalanceInt = sdk.ZeroInt()
	} else {
		
		usdtBalanceBigInt := k.erc20Keeper.BalanceOf(ctx, contracts2.ERC20MinterBurnerDecimalsContract.ABI, resp.TokenPair.GetERC20Contract(), evmAccount)

		usdtBalanceInt = sdk.NewIntFromBigInt(usdtBalanceBigInt)
		if usdtBalanceInt.IsNil() {
			usdtBalanceInt = sdk.ZeroInt()
		}
	}

	allBalances = append(allBalances, sdk.NewCoin(core.UsdtDenom, usdtBalanceInt))
	allBalances = allBalances.Sort()

	result, err := codec.MarshalJSONIndent(legacyQuerierCdc, allBalances)
	if err != nil {
		log.WithError(err).Error("MarshalJSONIndent queryMainTokenBalances Err")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return result, nil
}

func queryTokenPair(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainKeeper)
	var params erc20types.QueryTokenPairRequest
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	ctxs := sdk.WrapSDKContext(ctx)
	resp, err := k.erc20Keeper.TokenPair(ctxs, &params)
	if err != nil {
		return nil, err
	}
	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, resp.TokenPair)
	if err != nil {
		log.WithError(err).Error("MarshalJSONIndent TokenPair Err")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func queryContractCode(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainKeeper)
	var params evmtypes.QueryCodeRequest
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	addr := common.HexToAddress(params.Address)
	acct := k.evmKeeper.GetAccountWithoutBalance(ctx, addr)
	var code []byte
	if acct != nil && acct.IsContract() {
		code = k.evmKeeper.GetCode(ctx, common.BytesToHash(acct.CodeHash))
	}
	return code, nil
}


func queryParams(ctx sdk.Context, _ abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	params := k.GetParams(ctx)
	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
