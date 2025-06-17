package client

import (
	"encoding/hex"
	"encoding/json"
	"freemasonry.cc/blockchain/core/chainnet"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/ethereum/go-ethereum/common"
	evmhd "github.com/evmos/ethermint/crypto/hd"

	"github.com/tyler-smith/go-bip39"
)

func NewSecretKey() *SecretKey {
	return &SecretKey{}
}


type SecretKey struct {
}


func (k *SecretKey) CreateSeedWord() (string, error) {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return "", err
	}

	
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}
	return mnemonic, nil
}


func (k *SecretKey) CreateAccountFromSeed(mnemonic string) (*CosmosWallet, error) {
	keyringAlgos := keyring.SigningAlgoList{evmhd.EthSecp256k1, hd.Secp256k1}
	algo, err := keyring.NewSigningAlgoFromString(chainnet.Current.GetAlgo(), keyringAlgos)
	if err != nil {
		return nil, err
	}
	hdPath := hd.CreateHDPath(chainnet.Current.GetCoinType(), 0, 0).String()
	bip39Passphrase := ""
	derivedPriv, err := algo.Derive()(mnemonic, bip39Passphrase, hdPath)
	if err != nil {
		return nil, err
	}
	return k.CreateAccountFromPriv(hex.EncodeToString(derivedPriv))
}

func (k *SecretKey) CreateAccountFromPriv(priv string) (*CosmosWallet, error) {
	privKeyBytes, err := hex.DecodeString(priv)
	if err != nil {
		return nil, err
	}
	keyringAlgos := keyring.SigningAlgoList{evmhd.EthSecp256k1, hd.Secp256k1}
	algo, err := keyring.NewSigningAlgoFromString(chainnet.Current.GetAlgo(), keyringAlgos)
	if err != nil {
		return nil, err
	}
	privKey := algo.Generate()(privKeyBytes)

	bech32Addr, err := bech32.ConvertAndEncode(chainnet.Current.GetAddressPrefix(), privKey.PubKey().Address().Bytes())
	if err != nil {
		panic(err)
	}

	return &CosmosWallet{
		priv:       privKey,
		PrivateKey: priv,
		PublicKey:  hex.EncodeToString(privKey.PubKey().Bytes()),
		Address:    bech32Addr,
		EthAddress: common.BytesToAddress(privKey.PubKey().Address().Bytes()).String(),
	}, nil
}

func (k *SecretKey) Sign(addr *CosmosWallet, msg []byte) ([]byte, error) {
	return addr.priv.Sign(msg)
}

type CosmosWallet struct {
	Address    string        `json:"address"`
	EthAddress string        `json:"eth_address"`
	PublicKey  string        `json:"publickey"`
	PrivateKey string        `json:"privatekey"`
	priv       types.PrivKey `json:"priv"`
}

func (this *CosmosWallet) MarshalJson() []byte {
	data, _ := json.Marshal(this)
	return data
}
