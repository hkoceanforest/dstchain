package util

import (
	"encoding/hex"
	"freemasonry.cc/blockchain/core"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/types"
	"golang.org/x/crypto/sha3"
	"path/filepath"
)

func GenNodeAndValidatorKey(mnemonic, repoPath string) (p2p.ID, types.Address, error) {
	nodeKeyFile := filepath.Join(repoPath, "config", "node_key.json")
	pvKeyFile := filepath.Join(repoPath, "config", "priv_validator_key.json")
	pvStateFile := filepath.Join(repoPath, "data", "priv_validator_state.json")

	privKey := ed25519.GenPrivKeyFromSecret([]byte(mnemonic))
	filePV := privval.NewFilePV(privKey, pvKeyFile, pvStateFile)
	
	filePV.Key.Save()

	sha256 := sha3.New512()
	sha256.Write(privKey.Bytes())
	nodePrivkey := sha256.Sum(nil)
	nodeKey := &p2p.NodeKey{
		PrivKey: ed25519.PrivKey(nodePrivkey),
	}

	
	if err := nodeKey.SaveAs(nodeKeyFile); err != nil {
		return nodeKey.ID(), filePV.GetAddress(), err
	}

	return nodeKey.ID(), filePV.GetAddress(), nil
}

func GetAccountFromPub(pub string) (string, error) {
	pubKeyBytes, err := hex.DecodeString(pub)
	if err != nil {
		return "", err
	}

	pubkey := ethsecp256k1.PubKey{
		Key: pubKeyBytes,
	}

	address := sdk.AccAddress(pubkey.Address())
	return address.String(), nil
}

func GetPubKeyFromSign(sign, data []byte) (ethsecp256k1.PubKey, error) {
	dataUnSignC256 := crypto.Keccak256(data)
	secp256k1Pub, err := crypto.Ecrecover(dataUnSignC256, sign)
	if err != nil {
		return ethsecp256k1.PubKey{}, err
	}
	ecdsaPubkey, err := crypto.UnmarshalPubkey(secp256k1Pub)
	if err != nil {
		return ethsecp256k1.PubKey{}, err
	}
	comPub := ethsecp256k1.PubKey{Key: crypto.CompressPubkey(ecdsaPubkey)}

	ok := comPub.VerifySignature(data, sign)
	if !ok {
		return ethsecp256k1.PubKey{}, core.SignError
	}

	return comPub, nil
}
