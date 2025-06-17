package testenv

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"freemasonry.cc/blockchain/client"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/core/chainnet"
	"freemasonry.cc/blockchain/util"
	chattypes "freemasonry.cc/blockchain/x/chat/types"
	"freemasonry.cc/blockchain/x/dao/types"
	commtypes "freemasonry.cc/blockchain/x/gateway/types"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	types2 "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	types3 "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/crypto"
	evmhd "github.com/evmos/ethermint/crypto/hd"
	tmjson "github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/privval"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"time"
)



func (this *IntegrationTestSuite) NewContext() sdk.Context {
	this.Env.Height++
	return this.Env.App.BaseApp.NewContext(false, tmproto.Header{
		Height:  this.Env.Height,
		Time:    time.Now(),
		ChainID: chainnet.ChainIdDst,
	})
}

func (this *IntegrationTestSuite) printDecRate(dec sdk.Dec) string {
	return core.RemoveDecLastZero(dec.Mul(sdk.NewDec(100))) + "%"
}

func (this *IntegrationTestSuite) CreateNewAccount(ctx sdk.Context) types2.AccountI {
	priv := secp256k1.GenPrivKey()
	pk := priv.PubKey()
	newAccount := this.accountKep.NewAccountWithAddress(ctx, sdk.AccAddress(pk.Address()))
	this.accountKep.SetAccount(ctx, newAccount) 
	return newAccount
}

func (this *IntegrationTestSuite) CreateNewCluster(ctx sdk.Context, owner types2.AccountI, newMembers []types.Members, gatewayAddress, clusterId string) (types.DeviceCluster, error) {
	
	members := []types.Members{
		{
			MemberAddress: owner.GetAddress().String(),
			IndexNum:      "1111111",
			ChatAddress:   owner.GetAddress().String(),
		},
	}
	members = append(members, newMembers...)

	membersStrSlice := make([]string, 0)
	for _, member := range members {
		membersStrSlice = append(membersStrSlice, member.MemberAddress)
	}

	
	gatewaySign, onlineAmount, err := this.GetGatewaySign(1, membersStrSlice, false)
	this.T().Log(":", gatewaySign)
	this.T().Log(":", onlineAmount)

	clusterMsg := types.NewMsgCreateClusterAddMembers(
		owner.GetAddress().String(),
		gatewayAddress,
		clusterId,
		owner.GetAddress().String()+"chatAddress",
		clusterId+"clusterName",
		"1111111",
		gatewaySign,
		sdk.MustNewDecFromStr("0.5"),
		sdk.MustNewDecFromStr("0.5"),
		defaultBurnAmount,
		sdk.MustNewDecFromStr("0"),
		onlineAmount,
		members,
	)

	err = owner.SetSequence(1)
	this.Env.App.AccountKeeper.SetAccount(ctx, owner)
	this.Require().NoError(err)
	this.T().Log(":", owner.GetSequence())

	_, err = this.daoServ.CreateClusterAddMembers(ctx, clusterMsg)
	this.Require().NoError(err)
	this.T().Log("---------------------------------------------------------")
	
	cluster, err := this.daoKep.GetClusterByChatId(ctx, clusterMsg.ClusterId)
	this.Require().NoError(err)
	clusterInfo, err := json.Marshal(cluster)
	this.Require().NoError(err)
	this.T().Log(":", string(clusterInfo))

	return cluster, nil

}

func (this *IntegrationTestSuite) InitCluster(defaultBurnAmount sdk.Dec) (types.DeviceCluster, stakingtypes.ValidatorI, error) {
	app := this.Env.App
	env := this.Env
	env.Height = int64(1)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})
	
	amountDec := sdk.MustNewDecFromStr("100000000000000000000000")
	acc1Coins := sdk.NewCoins(
		sdk.NewCoin(core.BaseDenom, amountDec.TruncateInt()),
		sdk.NewCoin(core.GovDenom, amountDec.TruncateInt()),
	)
	
	this.Require().NoError(
		this.bankKep.SendCoins(ctx, genesisAccount, acc1.GetAddress(), acc1Coins),
	)
	this.T().Log(":", acc1.GetAddress().String())
	this.T().Log(":", app.BankKeeper.GetAllBalances(ctx, acc1.GetAddress()))
	
	validator, err := this.RegValidator(acc1.GetAddress(), pk1)
	this.Require().NoError(err)
	
	gatewayMsg := commtypes.NewMsgGatewayRegister(acc1.GetAddress().String(), "http://127.0.0.1:50327", "0", "a58981f1251b609192551a3ffd8acb04.a58", "12D3KooWKz5r8KenofDLZ4UciiX5eDZeT3spmJ256Ry7yjnKSpnc", "dst1e7y2f8mysu2u3m0agd8twr3fkhjdevh7z8m63g", []string{"1111111"})
	_, err = this.gatewayServ.GatewayRegister(ctx, gatewayMsg)
	this.Require().NoError(err)
	
	members := []types.Members{
		{
			MemberAddress: acc1.GetAddress().String(),
			IndexNum:      "1111111",
			ChatAddress:   acc1.GetAddress().String(),
		},
		
		
		
		
		
		
		
		
		
		
	}

	membersStrSlice := make([]string, 0)
	for _, member := range members {
		membersStrSlice = append(membersStrSlice, member.MemberAddress)
	}

	
	gatewaySign, onlineAmount, err := this.GetGatewaySign(1, membersStrSlice, false)
	this.T().Log(":", gatewaySign)
	this.T().Log(":", onlineAmount)

	clusterMsg := types.NewMsgCreateClusterAddMembers(
		acc1.GetAddress().String(),
		sdk.ValAddress(acc1.GetAddress()).String(),
		"!A7v7mxkI760rPyAB:1111111.nxn",
		"chatAddress",
		"clusterName",
		"1111111",
		gatewaySign,
		sdk.MustNewDecFromStr("0.5"),
		sdk.MustNewDecFromStr("0.5"),
		defaultBurnAmount,
		sdk.MustNewDecFromStr("0"),
		onlineAmount,
		members,
	)

	err = acc1.SetSequence(1)
	app.AccountKeeper.SetAccount(ctx, acc1)
	this.Require().NoError(err)
	this.T().Log(":", acc1.GetSequence())

	_, err = this.daoServ.CreateClusterAddMembers(ctx, clusterMsg)
	this.Require().NoError(err)
	this.T().Log("-------------------1--------------------------------------")
	
	cluster, err := this.daoKep.GetClusterByChatId(ctx, clusterMsg.ClusterId)
	
	this.Require().NoError(err)

	this.T().Log(fmt.Sprintf("【dst】,dst%s、%s、%s(:%s,:%s)、%sdao",
		this.printDecRate(this.daoKep.GetParams(ctx).BurnRegisterGateRatio),
		this.printDecRate(this.daoKep.GetParams(ctx).BurnCurrentGateRatio),
		this.printDecRate(this.daoKep.GetParams(ctx).DaoRewardPercent),
		this.printDecRate(cluster.ClusterDaoRatio),
		this.printDecRate(sdk.OneDec().Sub(cluster.ClusterDaoRatio)),
		this.printDecRate(this.daoKep.GetParams(ctx).Rate)),
	)

	this.T().Log(
		fmt.Sprintf("【dst】, * %sdaodao(:%s,%sdaodao)",
			this.printDecRate(this.daoKep.GetParams(ctx).BurnDaoPool),
			this.printDecRate(cluster.ClusterDaoRatio),
			this.printDecRate(sdk.OneDec().Sub(cluster.ClusterDaoRatio)),
		))

	this.T().Log(fmt.Sprintf("【dst】,dao,=dst * %s", this.printDecRate(this.daoKep.GetParams(ctx).ReceiveDaoRatio)))

	rate := sdk.OneDec().Sub(cluster.ClusterDeviceRatio)
	this.T().Log(fmt.Sprintf("【swapdst】,%sdao,dst:%s,:%sdstDao", this.printDecRate(this.daoKep.GetParams(ctx).BurnRewardFeeRate), this.printDecRate(rate), this.printDecRate(cluster.ClusterDeviceRatio)))

	this.T().Log(fmt.Sprintf("【dst】,:%s,:%s", this.printDecRate(cluster.ClusterSalaryRatio), this.printDecRate(sdk.OneDec().Sub(cluster.ClusterSalaryRatio))))

	this.T().Log("---------------------------------------------------")
	clusterByte, _ := json.Marshal(cluster)
	this.T().Log(":", string(clusterByte))
	return cluster, validator, nil
}

func (this *IntegrationTestSuite) CreateSmartValidator(accAddr sdk.AccAddress, coin sdk.Coin) (stakingtypes.ValidatorI, error) {
	env := this.Env
	app := env.App
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})
	rate, _ := sdk.NewDecFromStr("0.1")
	commission := stakingtypes.NewCommissionRates(rate, sdk.OneDec(), sdk.OneDec())
	description := stakingtypes.Description{
		Moniker:         "test",
		Identity:        "test",
		Website:         "test",
		SecurityContact: "test",
		Details:         "test",
	}
	validatorAddr := sdk.ValAddress(accAddr)
	msg, err := commtypes.NewMsgCreateSmartValidator(validatorAddr, "oK1vO9MU47i/IkPxI/oHdpp34759zQ1vFWw1VjwCJVo=", coin, description, commission, sdk.OneInt())
	if err != nil {
		return nil, err
	}
	uctx := sdk.WrapSDKContext(ctx)
	_, err = this.gatewayServ.CreateSmartValidator(uctx, msg)
	if err != nil {
		return nil, err
	}

	staking.EndBlocker(ctx, *this.stakingKep)

	
	env.Height += 1
	ctx = app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})

	
	staking.BeginBlocker(ctx, *this.stakingKep)
	validator, ok := this.stakingKep.GetValidator(ctx, validatorAddr)
	if !ok {
		return nil, errors.New("")
	}
	return validator, nil
}


func (this *IntegrationTestSuite) RegValidator(accAddr sdk.AccAddress, pk1 cryptotypes.PubKey) (stakingtypes.ValidatorI, error) {
	env := this.Env
	app := env.App
	
	ctx := env.Ctx
	uctx := sdk.WrapSDKContext(ctx)

	
	validatorAddr := sdk.ValAddress(accAddr)
	coinInt, _ := sdk.NewIntFromString("100000000000000000000")
	validatorDposCoin := sdk.NewCoin(core.GovDenom, coinInt)
	this.T().Log(":", validatorDposCoin)
	this.T().Log(":", validatorAddr.String())

	desc := stakingtypes.NewDescription("testname", "", "", "", "")
	comm := stakingtypes.NewCommissionRates(sdk.OneDec(), sdk.OneDec(), sdk.OneDec())

	msgCreateValidator, err := types3.NewMsgCreateValidator(validatorAddr, pk1, validatorDposCoin, desc, comm, sdk.NewInt(1))
	if err := msgCreateValidator.ValidateBasic(); err != nil {
		return nil, err
	}

	_, err = this.stakingServ.CreateValidator(uctx, msgCreateValidator)
	if err != nil {
		this.T().Fatal("：", err)
		return nil, err
	}
	staking.EndBlocker(ctx, *this.stakingKep)

	
	env.Height += 1
	ctx = app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})

	
	staking.BeginBlocker(ctx, *this.stakingKep)
	validator, ok := this.stakingKep.GetValidator(ctx, validatorAddr)
	if !ok {
		return nil, errors.New("")
	}
	return validator, nil
}


func (this *IntegrationTestSuite) RegGateway(accAddr sdk.AccAddress, delegation string, indexNumber []string) (gatewayAddr string, err error) {
	env := this.Env
	app := env.App
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})
	uctx := sdk.WrapSDKContext(ctx)

	gatewayRegisterMsg := commtypes.NewMsgGatewayRegister(accAddr.String(), "1", "192.168.0.117", delegation, "", "", indexNumber)

	if err = gatewayRegisterMsg.ValidateBasic(); err != nil {
		return
	}

	_, err = this.gatewayServ.GatewayRegister(uctx, gatewayRegisterMsg)
	if err != nil {
		return
	}
	gatewayAddr = sdk.ValAddress(accAddr).String()
	return
}


func (this *IntegrationTestSuite) GatewayDistributeNumber(accAddr sdk.AccAddress, validatorAddr sdk.ValAddress, numbers []string) (err error) {
	env := this.Env
	app := env.App
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})
	uctx := sdk.WrapSDKContext(ctx)

	msgGatewayIndexNum := commtypes.NewMsgGatewayIndexNum(accAddr.String(), validatorAddr.String(), numbers)
	if err = msgGatewayIndexNum.ValidateBasic(); err != nil {
		return err
	}

	
	_, err = this.gatewayServ.GatewayIndexNum(uctx, msgGatewayIndexNum)

	return err
}

func (this *IntegrationTestSuite) GatewayDelegation(accAddr sdk.AccAddress, validatorAddr sdk.ValAddress, coin sdk.Coin) (err error) {
	env := this.Env
	app := env.App
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})
	uctx := sdk.WrapSDKContext(ctx)

	msgDelegation := stakingtypes.NewMsgDelegate(accAddr, validatorAddr, coin)

	
	_, err = this.stakingServ.Delegate(uctx, msgDelegation)
	if err != nil {
		return err
	}
	env.Height += 1
	ctx = app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})
	
	staking.BeginBlocker(ctx, *this.stakingKep)
	staking.EndBlocker(ctx, *this.stakingKep)
	return nil
}

func (this *IntegrationTestSuite) GatewayUndelegation(accAddr sdk.AccAddress, validatorAddr sdk.ValAddress, coin sdk.Coin, numbers []string) (err error) {
	env := this.Env
	app := env.App
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})
	uctx := sdk.WrapSDKContext(ctx)

	msgUndelegation := commtypes.NewMsgGatewayUndelegation(accAddr.String(), validatorAddr.String(), coin, numbers)

	
	_, err = this.gatewayServ.GatewayUndelegate(uctx, msgUndelegation)
	if err != nil {
		return err
	}
	env.Height += 1
	ctx = app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})
	
	staking.BeginBlocker(ctx, *this.stakingKep)
	staking.EndBlocker(ctx, *this.stakingKep)
	return nil
}


func (this *IntegrationTestSuite) BurnGetMobile(accAddr sdk.AccAddress, mobilePrefix string) (err error) {
	env := this.Env
	app := env.App
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})
	uctx := sdk.WrapSDKContext(ctx)

	msg := chattypes.NewMsgBurnGetMobile(accAddr.String(), mobilePrefix, "", "")

	if err = msg.ValidateBasic(); err != nil {
		return err
	}
	_, err = this.chatServ.BurnGetMobile(uctx, msg)
	if err != nil {
		this.T().Log("")
		this.T().Fatal(err)
		return err
	}
	return nil
}


func (this *IntegrationTestSuite) MobileTransfer(fromAddr sdk.AccAddress, toAddr sdk.AccAddress, mobileNumber string) (err error) {
	env := this.Env
	app := env.App
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})
	uctx := sdk.WrapSDKContext(ctx)

	msg := chattypes.NewMsgMobileTransfer(fromAddr.String(), toAddr.String(), mobileNumber)

	if err = msg.ValidateBasic(); err != nil {
		return err
	}
	_, err = this.chatServ.MobileTransfer(uctx, msg)
	if err != nil {
		this.T().Log("")
		this.T().Fatal(err)
		return err
	}
	return nil
}


func (this *IntegrationTestSuite) SendCoinsFromModuleToAccount(accAddr sdk.AccAddress, coins sdk.Coins) error {
	app := this.Env.App
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: this.Env.Height})
	return this.bankKep.SendCoinsFromModuleToAccount(ctx, moduleAccountName, accAddr, coins)
}


func (this *IntegrationTestSuite) CreateCluster(num int) ([]types.DeviceCluster, stakingtypes.ValidatorI, error) {
	app := this.Env.App
	env := this.Env
	env.Height = int64(1)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})
	
	amountDec := sdk.MustNewDecFromStr("200000000000000000000000000")
	nxnAmountDec := sdk.MustNewDecFromStr("100000000000000000000000")
	acc1Coins := sdk.NewCoins(
		sdk.NewCoin(core.BaseDenom, amountDec.TruncateInt()),
		sdk.NewCoin(core.GovDenom, nxnAmountDec.TruncateInt()),
	)
	accCoins := sdk.NewCoins(
		sdk.NewCoin(core.BaseDenom, amountDec.TruncateInt()),
	)
	
	this.Require().NoError(
		this.bankKep.SendCoins(ctx, genesisAccount, acc1.GetAddress(), acc1Coins),
	)
	this.T().Log(":", acc1.GetAddress().String())
	this.T().Log(":", app.BankKeeper.GetAllBalances(ctx, acc1.GetAddress()))
	
	validator, err := this.RegValidator(acc1.GetAddress(), pk1)
	this.Require().NoError(err)
	machineAddress := "dst145kjpzns5awv7lwn8vjl0qkuhqy0kcdlp43h4r"
	
	gatewayMsg := commtypes.NewMsgGatewayRegister(acc1.GetAddress().String(), "http://127.0.0.1:50327", "0", "a58981f1251b609192551a3ffd8acb04.a58", "12D3KooWQy4Rycyj5xVZDjYCMxH3ZEDThPboKuWP8TatrV6f6kn3", machineAddress, []string{"1111111"})
	_, err = this.gatewayServ.GatewayRegister(ctx, gatewayMsg)
	this.Require().NoError(err)
	clusters := make([]types.DeviceCluster, 0)
	txClient := client.NewTxClient()
	accountClient := client.NewAccountClient(&txClient)
	for i := 0; i < num; i++ {
		seed, _ := accountClient.CreateSeedWord()
		toWallet, _ := accountClient.CreateAccountFromSeed(seed)
		from, _ := sdk.AccAddressFromBech32(toWallet.Address)
		this.Require().NoError(
			this.bankKep.SendCoins(ctx, genesisAccount, from, accCoins),
		)
		
		randStr, _ := util.GenerateSecureRandomString(16)
		
		member := types.Members{MemberAddress: toWallet.Address, IndexNum: "1111111", ChatAddress: ""}
		members := []types.Members{member}
		clusterMsg := types.NewMsgCreateClusterAddMembers(machineAddress, validator.GetOperator().String(), "!"+randStr+":1111111.nxn", toWallet.Address, toWallet.Address, toWallet.Address, "", sdk.NewDecWithPrec(10, 2), sdk.NewDecWithPrec(90, 2), defaultBurnAmount, sdk.NewDec(0), 0, members)
		_, err = this.daoServ.CreateClusterAddMembers(ctx, clusterMsg)
		this.Require().NoError(err)
		
		cluster, err := this.daoKep.GetClusterByChatId(ctx, clusterMsg.ClusterId)
		this.Require().NoError(err)
		clusters = append(clusters, cluster)
	}
	return clusters, validator, nil
}

func (this *IntegrationTestSuite) GetGatewaySign(seq int64, members []string, memberAdd bool) (string, int64, error) {
	
	types.SortSliceMembers(members)
	membersBytes, err := json.Marshal(members)
	if err != nil {
		return "", 0, err
	}

	
	onlineCount := int64(1)

	signData := types.GatewaySign{
		Members:      string(membersBytes),
		OnlineAmount: onlineCount,
		Seq:          seq - 1,
		MemberAdd:    memberAdd,
	}
	
	
	signDataBytes, err := json.Marshal(signData)
	if err != nil {
		return "", 0, err
	}

	signOriginalData := crypto.Keccak256(signDataBytes)

	
	
	pvKey := privval.FilePVKey{}
	err = tmjson.Unmarshal([]byte(priv_dalidator_key), &pvKey)
	if err != nil {
		return "", 0, err
	}

	privatekey := hex.EncodeToString(pvKey.PrivKey.Bytes())

	gatewaySign, err := Sign(privatekey, signOriginalData)
	if err != nil {
		return "", 0, err
	}

	return gatewaySign, onlineCount, nil
}


func Sign(privateKey string, msg []byte) (string, error) {
	privKeyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return "", err
	}
	keyringAlgos := keyring.SigningAlgoList{evmhd.EthSecp256k1, hd.Secp256k1}
	algo, err := keyring.NewSigningAlgoFromString(chainnet.Current.GetAlgo(), keyringAlgos)
	if err != nil {
		return "", err
	}
	privKey := algo.Generate()(privKeyBytes)
	priv := privKey
	signedTxBytes, err := priv.Sign(msg)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(signedTxBytes), nil
}
