package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"encoding/json"
	"fmt"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/x/contract/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingKeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/evmos/ethermint/server/config"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	erc20keeper "github.com/evmos/evmos/v10/x/erc20/keeper"
	erc20types "github.com/evmos/evmos/v10/x/erc20/types"
	"github.com/tendermint/tendermint/libs/log"
	"math/big"
	"strings"
)

type Keeper struct {
	storeKey   storetypes.StoreKey
	cdc        codec.BinaryCodec
	paramstore paramtypes.Subspace

	stakingKeeper *stakingKeeper.Keeper
	accountKeeper types.AccountKeeper
	BankKeeper    types.BankKeeper
	evmKeeper     types.EVMKeeper
	chatKeeper    types.ChatKeeper
	gatewayKeeper types.GatewayKeeper
	erc20Keeper   erc20keeper.Keeper
}

func NewKeeper(
	storeKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
	ps paramtypes.Subspace,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	stakingKeeper *stakingKeeper.Keeper,
	evmKeeper types.EVMKeeper,
	chatKeeper types.ChatKeeper,
	gatewayKeeper types.GatewayKeeper,
	erc20Keeper erc20keeper.Keeper,
) Keeper {
	
	
	
	

	return Keeper{
		storeKey:      storeKey,
		cdc:           cdc,
		paramstore:    ps,
		accountKeeper: ak,
		BankKeeper:    bk,
		stakingKeeper: stakingKeeper,
		evmKeeper:     evmKeeper,
		chatKeeper:    chatKeeper,
		gatewayKeeper: gatewayKeeper,
		erc20Keeper:   erc20Keeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) KVHelper(ctx sdk.Context) StoreHelper {
	store := ctx.KVStore(k.storeKey)
	return StoreHelper{
		store,
	}
}


func (k Keeper) IsContract(ctx sdk.Context, contract string) bool {
	address := k.evmKeeper.GetAccountWithoutBalance(ctx, common.HexToAddress(contract))
	if address == nil {
		return false
	}
	return address.IsContract()
}

func (k Keeper) IsERC20Enabled(ctx sdk.Context) bool {
	return k.erc20Keeper.IsERC20Enabled(ctx)
}

func (k Keeper) GetTokenPairID(ctx sdk.Context, token string) []byte {
	return k.erc20Keeper.GetTokenPairID(ctx, token)
}
func (k Keeper) GetTokenPair(ctx sdk.Context, id []byte) (erc20types.TokenPair, bool) {
	return k.erc20Keeper.GetTokenPair(ctx, id)
}


func (k Keeper) QueryContractIsExist(ctx sdk.Context, address common.Address) bool {
	acct := k.evmKeeper.GetAccountWithoutBalance(ctx, address)
	var code []byte
	if acct != nil && acct.IsContract() {
		code = k.evmKeeper.GetCode(ctx, common.BytesToHash(acct.CodeHash))
	}
	if len(code) == 0 {
		return false
	}
	return true
}

func (k Keeper) CallEVM(
	ctx sdk.Context,
	abi abi.ABI,
	from, contract common.Address,
	commit bool,
	method string,
	args ...interface{},
) (*evmtypes.MsgEthereumTxResponse, error) {
	data, err := abi.Pack(method, args...)
	if err != nil {
		return nil, sdkerrors.Wrap(
			core.ErrABIPack,
			sdkerrors.Wrap(err, "failed to create transaction data").Error(),
		)
	}

	resp, err := k.CallEVMWithData(ctx, from, &contract, data, commit)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "contract call failed: method '%s', contract '%s'", method, contract)
	}
	return resp, nil
}

func (k Keeper) CallEVMWithData(
	ctx sdk.Context,
	from common.Address,
	contract *common.Address,
	data []byte,
	commit bool,
) (*evmtypes.MsgEthereumTxResponse, error) {
	nonce, err := k.accountKeeper.GetSequence(ctx, from.Bytes())
	if err != nil && !strings.Contains(err.Error(), "does not exist") {
		return nil, err
	}
	gasCap := config.DefaultGasCap
	if commit {
		args, err := json.Marshal(evmtypes.TransactionArgs{
			From: &from,
			To:   contract,
			Data: (*hexutil.Bytes)(&data),
		})
		if err != nil {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, "failed to marshal tx args: %s", err.Error())
		}

		gasRes, err := k.evmKeeper.EstimateGas(sdk.WrapSDKContext(ctx), &evmtypes.EthCallRequest{
			Args:   args,
			GasCap: config.DefaultGasCap,
		})
		if err != nil {
			return nil, err
		}
		gasCap = gasRes.Gas
	}
	msg := ethtypes.NewMessage(
		from,
		contract,
		nonce,
		big.NewInt(0), 
		gasCap,        
		big.NewInt(0), 
		big.NewInt(0), 
		big.NewInt(0), 
		data,
		ethtypes.AccessList{}, 
		!commit,               
	)
	res, err := k.evmKeeper.ApplyMessage(ctx, msg, evmtypes.NewNoOpTracer(), commit)
	if err != nil {
		return nil, err
	}

	if res.Failed() {
		return nil, sdkerrors.Wrap(evmtypes.ErrVMExecution, res.VmError)
	}

	return res, nil
}

func (k Keeper) CallEVMWithValue(
	ctx sdk.Context,
	abi abi.ABI,
	from, contract common.Address,
	commit bool,
	method string,
	value *hexutil.Big,
	args ...interface{},
) (*evmtypes.MsgEthereumTxResponse, error) {
	data, err := abi.Pack(method, args...)
	if err != nil {
		return nil, sdkerrors.Wrap(
			core.ErrABIPack,
			sdkerrors.Wrap(err, "failed to create transaction data").Error(),
		)
	}

	fmt.Printf("Calldata: 0x%x\n", data)

	resp, err := k.CallEVMWithDataValue(ctx, from, &contract, data, commit, value)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "contract call failed: method '%s', contract '%s'", method, contract)
	}
	return resp, nil
}

func (k Keeper) CallEVMWithDataValue(
	ctx sdk.Context,
	from common.Address,
	contract *common.Address,
	data []byte,
	commit bool,
	value *hexutil.Big,
) (*evmtypes.MsgEthereumTxResponse, error) {
	nonce, err := k.accountKeeper.GetSequence(ctx, from.Bytes())
	if err != nil && !strings.Contains(err.Error(), "does not exist") {
		return nil, err
	}
	gasCap := config.DefaultGasCap
	if commit {
		args, err := json.Marshal(evmtypes.TransactionArgs{
			From:  &from,
			To:    contract,
			Data:  (*hexutil.Bytes)(&data),
			Value: value,
		})
		if err != nil {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, "failed to marshal tx args: %s", err.Error())
		}

		gasRes, err := k.evmKeeper.EstimateGas(sdk.WrapSDKContext(ctx), &evmtypes.EthCallRequest{
			Args:   args,
			GasCap: config.DefaultGasCap,
		})
		if err != nil {
			return nil, err
		}
		gasCap = gasRes.Gas
	}
	msg := ethtypes.NewMessage(
		from,
		contract,
		nonce,
		value.ToInt(), 
		gasCap,        
		big.NewInt(0), 
		big.NewInt(0), 
		big.NewInt(0), 
		data,
		ethtypes.AccessList{}, 
		!commit,               
	)
	res, err := k.evmKeeper.ApplyMessage(ctx, msg, evmtypes.NewNoOpTracer(), commit)
	if err != nil {
		return nil, err
	}

	if res.Failed() {
		return nil, sdkerrors.Wrap(evmtypes.ErrVMExecution, res.VmError)
	}

	return res, nil
}


func (k Keeper) SetTokenFactoryContractAddress(ctx sdk.Context, contract string) error {
	store := k.KVHelper(ctx)
	err := store.Set(types.KeyTokenFactoryContractAddress, contract)
	if err != nil {
		return err
	}
	return nil
}


func (k Keeper) SetRedPacketContractAddress(ctx sdk.Context, contract string) error {
	store := k.KVHelper(ctx)
	err := store.Set(types.KeyRedPacketContractAddress, contract)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) GetRedPacketContractAddress(ctx sdk.Context) string {
	store := k.KVHelper(ctx)

	contractAddress := string(store.Get(types.KeyRedPacketContractAddress))

	if contractAddress == "" {
		return strings.ToLower("0x570d45E28C4B3AeC78F25533981Ae89eB5a564dC")
	}
	return strings.ToLower(contractAddress)
}

func (k Keeper) GetTokenFactoryContractAddress(ctx sdk.Context) string {
	store := k.KVHelper(ctx)

	contractAddress := string(store.Get(types.KeyTokenFactoryContractAddress))

	if contractAddress == "" {
		return strings.ToLower("0xC09E37c186dc89aff7244AFBd63190Ed1e04E4AA")
	}
	return strings.ToLower(contractAddress)
}

func (k Keeper) RegisterERC20(ctx sdk.Context, contractAddress, daemon string, owner int32) (*erc20types.TokenPair, error) {
	contract := common.HexToAddress(contractAddress)
	
	if k.erc20Keeper.IsERC20Registered(ctx, contract) {
		return nil, errorsmod.Wrapf(
			erc20types.ErrTokenPairAlreadyExists, "token ERC20 contract already registered: %s", contract.String(),
		)
	}

	metadata, err := k.CreateCoinMetadata(ctx, contract, daemon)
	if err != nil {
		return nil, errorsmod.Wrap(
			err, "failed to create wrapped coin denom metadata for ERC20",
		)
	}
	var contractOwner erc20types.Owner
	switch owner {
	case 1:
		contractOwner = erc20types.OWNER_MODULE
	case 2:
		contractOwner = erc20types.OWNER_EXTERNAL
	default:
		contractOwner = erc20types.OWNER_EXTERNAL
	}
	pair := erc20types.NewTokenPair(contract, metadata.Name, true, contractOwner)
	k.erc20Keeper.SetTokenPair(ctx, pair)
	k.erc20Keeper.SetDenomMap(ctx, pair.Denom, pair.GetID())
	k.erc20Keeper.SetERC20Map(ctx, common.HexToAddress(pair.Erc20Address), pair.GetID())
	return &pair, nil
}

func (k Keeper) RegisterCoin(ctx sdk.Context, coinMetadata banktypes.Metadata) (*erc20types.TokenPair, error) {
	tokenPair, err := k.erc20Keeper.RegisterCoin(ctx, coinMetadata)
	if err != nil {
		return nil, err
	}
	return tokenPair, nil
}

func (k Keeper) splitCoin(amount sdk.Coins) (sdk.Coins, sdk.Coin) {
	usdtCoin := sdk.NewCoin(core.UsdtDenom, sdk.ZeroInt())
	newCoins := sdk.Coins{}
	for _, coin := range amount {
		if coin.Denom == core.UsdtDenom {
			usdtCoin = coin
			continue
		}
		newCoins = append(newCoins, coin)
	}
	return newCoins, usdtCoin
}

func (k Keeper) CreateCoinMetadata(
	ctx sdk.Context,
	contract common.Address,
	daemon string,
) (*banktypes.Metadata, error) {
	strContract := contract.String()
	erc20Data, err := k.erc20Keeper.QueryERC20(ctx, contract)
	if err != nil {
		return nil, err
	}

	
	_, found := k.BankKeeper.GetDenomMetaData(ctx, erc20types.CreateDenom(strContract))
	if found {
		return nil, errorsmod.Wrap(
			erc20types.ErrInternalTokenPair, "denom metadata already registered",
		)
	}

	if k.erc20Keeper.IsDenomRegistered(ctx, erc20types.CreateDenom(strContract)) {
		return nil, errorsmod.Wrapf(
			erc20types.ErrInternalTokenPair, "coin denomination already registered: %s", erc20Data.Name,
		)
	}

	
	base := erc20types.CreateDenom(strContract)
	if daemon != "" {
		base = daemon
	}
	
	
	
	metadata := banktypes.Metadata{
		Description: erc20types.CreateDenomDescription(strContract),
		Base:        base,
		
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    base,
				Exponent: 0,
			},
		},
		Name:    base,
		Symbol:  erc20Data.Symbol,
		Display: base,
	}

	
	if erc20Data.Decimals > 0 {
		nameSanitized := erc20types.SanitizeERC20Name(erc20Data.Name)
		metadata.DenomUnits = append(
			metadata.DenomUnits,
			&banktypes.DenomUnit{
				Denom:    nameSanitized,
				Exponent: uint32(erc20Data.Decimals),
			},
		)
		metadata.Display = nameSanitized
	}

	if err := metadata.Validate(); err != nil {
		return nil, errorsmod.Wrapf(
			err, "ERC20 token data is invalid for contract %s", strContract,
		)
	}

	k.BankKeeper.SetDenomMetaData(ctx, metadata)

	return &metadata, nil
}
