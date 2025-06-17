package client

import (
	"encoding/hex"
	"errors"
	"fmt"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/core/chainnet"
	"freemasonry.cc/blockchain/util"
	"freemasonry.cc/blockchain/x/contract/types"
	daoTypes "freemasonry.cc/blockchain/x/dao/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	types2 "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	types3 "github.com/evmos/ethermint/types"
	"github.com/sirupsen/logrus"
	ttypes "github.com/tendermint/tendermint/types"
	"strings"
)

type AccountClient struct {
	TxClient  *TxClient
	key       *SecretKey
	logPrefix string
}

type Account struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Address  string `json:"address"`
	Pubkey   string `json:"pubkey"`
	Mnemonic string `'json:"mnemonic"`
}

func (this *Account) Print() {
	fmt.Printf("Name:\t %s \n", this.Name)
	fmt.Printf("Address:\t %s \n", this.Address)
	fmt.Printf("Type:\t\t %s \n", this.Type)
	fmt.Printf("Pubkey:\t\t %s \n", this.Pubkey)
	fmt.Printf("Menmonic:\t\t %s \n", this.Mnemonic)
}

type AccountList struct {
	Accounts []Account
}

func (this *AccountClient) Transfer(data core.TransferData, privateKey string) (tx ttypes.Tx, resp *core.BroadcastTxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	from, err := sdk.AccAddressFromBech32(data.FromAddress)
	if err != nil {
		return tx, nil, err
	}
	to, err := sdk.AccAddressFromBech32(data.ToAddress)
	if err != nil {
		return tx, nil, err
	}
	
	msg := banktypes.NewMsgSend(from, to, data.Coins)
	if err != nil {
		log.Error("NewMsgTransfer")
		return
	}
	var result *core.BaseResponse
	
	tx, result, err = this.TxClient.SignAndSendMsg(data.FromAddress, privateKey, data.Fee, data.Memo, msg)
	if err != nil {
		return
	}
	resp = new(core.BroadcastTxResponse)
	
	if result.Status == 1 {
		dataByte, err1 := util.Json.Marshal(result.Data)
		if err1 != nil {
			err = err1
			return
		}
		err = util.Json.Unmarshal(dataByte, resp)
		if err != nil {
			return
		}
		return tx, resp, nil
	} else {
		
		resp.TxHash = hex.EncodeToString(tx.Hash())
		return tx, resp, errors.New(result.Info)
	}
}


func (this *AccountClient) GetAllAccounts() (accounts []string, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	reponseStr, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/auth/accounts", []byte{})
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return
	}
	err = clientCtx.LegacyAmino.UnmarshalJSON(reponseStr, &accounts)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON2")
		return
	}
	return
}

func (this *AccountClient) FindAccountBalances(accountAddr string, height string) (coins sdk.Coins, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient).WithFields(logrus.Fields{"acc": accountAddr})
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	req := banktypes.QueryAllBalancesRequest{Address: accountAddr}
	reqBytes, _ := clientCtx.LegacyAmino.MarshalJSON(req)

	reponseStr, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/bank/all_balances", reqBytes)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return
	}
	err = clientCtx.LegacyAmino.UnmarshalJSON(reponseStr, &coins)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON2")
		return
	}
	return
}

func (this *AccountClient) FindAccountBalance(accountAddr string, denom, height string) (coin sdk.Coin, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient).WithFields(logrus.Fields{"acc": accountAddr, "denom": denom})
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	req := banktypes.QueryBalanceRequest{Address: accountAddr, Denom: denom}

	reqBytes, _ := clientCtx.LegacyAmino.MarshalJSON(req)

	reponseStr, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/bank/balance", reqBytes)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return
	}
	err = clientCtx.LegacyAmino.UnmarshalJSON(reponseStr, &coin)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON2")
		return
	}
	return
}

func (this *AccountClient) FindBalanceByRpc(accountAddr string, denom string) (realCoins core.RealCoin, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient).WithFields(logrus.Fields{"acc": accountAddr, "denom": denom})
	clientCtx := chainnet.ChainNetDst.GetClientCtx()

	req := banktypes.QueryBalanceRequest{Address: accountAddr, Denom: denom}

	reqBytes, _ := clientCtx.LegacyAmino.MarshalJSON(req)

	reponseStr, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/bank/balance", reqBytes)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return
	}
	var coin sdk.Coin
	err = clientCtx.LegacyAmino.UnmarshalJSON(reponseStr, &coin)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON2")
		return
	}
	realCoins = core.MustLedgerCoin2RealCoin(coin)
	return
}

func (this *AccountClient) CreateAccountFromPriv(priv string) (*CosmosWallet, error) {
	return this.key.CreateAccountFromPriv(priv)
}



func (this *AccountClient) CreateAccountFromSeed(seed string) (acc *CosmosWallet, err error) {
	
	seed = strings.TrimPrefix(seed, " ")
	seed = strings.TrimSuffix(seed, " ")
	seeds := strings.Split(seed, " ")

	filteredSeed := []string{} 
	for _, v := range seeds {
		if v != "" {
			filteredSeed = append(filteredSeed, v)
		}
	}
	return this.key.CreateAccountFromSeed(strings.Join(filteredSeed, " "))
}


func (this *AccountClient) CreateSeedWord() (mnemonic string, err error) {
	return this.key.CreateSeedWord()
}

func (this *AccountClient) FindMainTokenBalances(address string) (coins sdk.Coins, err error) {
	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return
	}
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	
	bz, err := clientCtx.LegacyAmino.MarshalJSON(addr)
	if err != nil {
		return
	}
	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryMainTokenBalances)
	tokenbalanceinfo, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, route, bz)
	if err != nil {
		return
	}

	if tokenbalanceinfo == nil {
		return nil, core.ErrQueryTokens
	}

	var res sdk.Coins
	err = clientCtx.LegacyAmino.UnmarshalJSON(tokenbalanceinfo, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}


func (this *AccountClient) GetAccount(addr string) (acc types3.EthAccount, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	params := types2.QueryAccountRequest{
		Address: addr,
	}

	paramsData, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return
	}

	reponseStr, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/auth/account", paramsData)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return
	}
	var account types3.EthAccount
	err = clientCtx.LegacyAmino.UnmarshalJSON(reponseStr, &account)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON2")
		return
	}
	return account, nil
}

func (this *AccountClient) GetNextAccountNumber() (accNum uint64, err error) {
	logs := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/dao/"+daoTypes.QueryNextAccountNumber, nil)
	if err != nil {
		logs.WithError(err).Error("QueryWithData")
		return accNum, err
	}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &accNum)
	if err != nil {
		logs.WithError(err).Error("UnmarshalJSON")
		return accNum, err
	}
	return accNum, nil

}
