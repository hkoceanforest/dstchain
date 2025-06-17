package cli

import (
	"encoding/hex"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	evmhd "github.com/evmos/ethermint/crypto/hd"
	"io/ioutil"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/codec"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func ParseMetadata(cdc codec.JSONCodec, metadataFile string) (banktypes.Metadata, error) {
	metadata := banktypes.Metadata{}

	contents, err := ioutil.ReadFile(filepath.Clean(metadataFile))
	if err != nil {
		return metadata, err
	}

	if err = cdc.UnmarshalJSON(contents, &metadata); err != nil {
		return metadata, err
	}

	return metadata, nil
}

func createAccount(priv string) (string, error) {
	privKeyBytes, err := hex.DecodeString(priv)
	if err != nil {
		return "", err
	}
	keyringAlgos := keyring.SigningAlgoList{evmhd.EthSecp256k1, hd.Secp256k1}
	algo, err := keyring.NewSigningAlgoFromString(ethsecp256k1.KeyType, keyringAlgos)
	if err != nil {
		return "", err
	}
	privKey := algo.Generate()(privKeyBytes)

	bech32Addr, err := bech32.ConvertAndEncode(sdk.Bech32MainPrefix, privKey.PubKey().Address().Bytes())
	if err != nil {
		panic(err)
	}
	return bech32Addr, nil
}
