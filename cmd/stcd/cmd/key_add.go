package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	cryptokeyring "github.com/cosmos/cosmos-sdk/crypto/keyring"
	etherminthd "github.com/evmos/ethermint/crypto/hd"
	"io"
	"sort"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bip39 "github.com/cosmos/go-bip39"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

const (
	flagInteractive       = "interactive"
	flagRecover           = "recover"
	flagNoBackup          = "no-backup"
	flagCoinType          = "coin-type"
	flagAccount           = "account"
	flagIndex             = "index"
	flagMultisig          = "multisig"
	flagMultiSigThreshold = "multisig-threshold"
	flagNoSort            = "nosort"
	flagHDPath            = "hd-path"

	mnemonicEntropySize = 256
)


func KeyAddCmd(ctx client.Context, cmd *cobra.Command, args []string, inBuf *bufio.Reader) error {
	var (
		algo keyring.SignatureAlgo
		err  error
	)

	name := args[0]

	interactive, _ := cmd.Flags().GetBool(flagInteractive)
	noBackup, _ := cmd.Flags().GetBool(flagNoBackup)
	useLedger, _ := cmd.Flags().GetBool(flags.FlagUseLedger)
	algoStr, _ := cmd.Flags().GetString(flags.FlagKeyAlgorithm)

	showMnemonic := !noBackup
	kb := ctx.Keyring
	outputFormat := ctx.OutputFormat

	keyringAlgos, ledgerAlgos := kb.SupportedAlgorithms()

	
	
	if useLedger {
		algo, err = keyring.NewSigningAlgoFromString(algoStr, ledgerAlgos)
	} else {
		algo, err = keyring.NewSigningAlgoFromString(algoStr, keyringAlgos)
	}

	if err != nil {
		return err
	}

	if dryRun, _ := cmd.Flags().GetBool(flags.FlagDryRun); dryRun {
		
		kb = keyring.NewInMemory(ctx.Codec, etherminthd.EthSecp256k1Option())
	} else {
		_, err = kb.Key(name)
		if err == nil {
			
			response, err2 := input.GetConfirmation(fmt.Sprintf("override the existing name %s", name), inBuf, cmd.ErrOrStderr())
			if err2 != nil {
				return err2
			}

			if !response {
				return errors.New("aborted")
			}

			err2 = kb.Delete(name)
			if err2 != nil {
				return err2
			}
		}

		multisigKeys, _ := cmd.Flags().GetStringSlice(flagMultisig)
		if len(multisigKeys) != 0 {
			pks := make([]cryptotypes.PubKey, len(multisigKeys))
			multisigThreshold, _ := cmd.Flags().GetInt(flagMultiSigThreshold)
			if err := validateMultisigThreshold(multisigThreshold, len(multisigKeys)); err != nil {
				return err
			}

			for i, keyname := range multisigKeys {
				k, err := kb.Key(keyname)
				if err != nil {
					return err
				}

				key, err := k.GetPubKey()
				if err != nil {
					return err
				}
				pks[i] = key
			}

			if noSort, _ := cmd.Flags().GetBool(flagNoSort); !noSort {
				sort.Slice(pks, func(i, j int) bool {
					return bytes.Compare(pks[i].Address(), pks[j].Address()) < 0
				})
			}

			pk := multisig.NewLegacyAminoPubKey(multisigThreshold, pks)
			k, err := kb.SaveMultisig(name, pk)
			if err != nil {
				return err
			}

			return printCreate(cmd, k, false, "", outputFormat)
		}
	}

	pubKey, _ := cmd.Flags().GetString(keys.FlagPublicKey)
	if pubKey != "" {
		var pk cryptotypes.PubKey
		if err = ctx.Codec.UnmarshalInterfaceJSON([]byte(pubKey), &pk); err != nil {
			return err
		}

		k, err := kb.SaveOfflineKey(name, pk)
		if err != nil {
			return err
		}

		return printCreate(cmd, k, false, "", outputFormat)
	}

	coinType, _ := cmd.Flags().GetUint32(flagCoinType)
	account, _ := cmd.Flags().GetUint32(flagAccount)
	index, _ := cmd.Flags().GetUint32(flagIndex)
	hdPath, _ := cmd.Flags().GetString(flagHDPath)

	if len(hdPath) == 0 {
		hdPath = hd.CreateHDPath(coinType, account, index).String()
	} else if useLedger {
		return errors.New("cannot set custom bip32 path with ledger")
	}

	
	if useLedger {
		bech32PrefixAccAddr := sdk.GetConfig().GetBech32AccountAddrPrefix()

		
		k, err := kb.SaveLedgerKey(name, algo, bech32PrefixAccAddr, coinType, account, index)
		if err != nil {
			return err
		}

		return printCreate(cmd, k, false, "", outputFormat)
	}

	
	var mnemonic, bip39Passphrase string

	recover, _ := cmd.Flags().GetBool(flagRecover)
	if recover {
		mnemonic, err = input.GetString("Enter your bip39 mnemonic", inBuf)
		if err != nil {
			return err
		}

		if !bip39.IsMnemonicValid(mnemonic) {
			return errors.New("invalid mnemonic")
		}
	} else if interactive {
		mnemonic, err = input.GetString("Enter your bip39 mnemonic, or hit enter to generate one.", inBuf)
		if err != nil {
			return err
		}

		if !bip39.IsMnemonicValid(mnemonic) && mnemonic != "" {
			return errors.New("invalid mnemonic")
		}
	}

	if len(mnemonic) == 0 {
		
		entropySeed, err := bip39.NewEntropy(128)
		if err != nil {
			return err
		}

		mnemonic, err = bip39.NewMnemonic(entropySeed)
		if err != nil {
			return err
		}
	}

	
	if interactive {
		bip39Passphrase, err = input.GetString(
			"Enter your bip39 passphrase. This is combined with the mnemonic to derive the seed. "+
				"Most users should just hit enter to use the default, \"\"", inBuf)
		if err != nil {
			return err
		}

		
		if len(bip39Passphrase) != 0 {
			p2, err := input.GetString("Repeat the passphrase:", inBuf)
			if err != nil {
				return err
			}

			if bip39Passphrase != p2 {
				return errors.New("passphrases don't match")
			}
		}
	}

	k, err := kb.NewAccount(name, mnemonic, bip39Passphrase, hdPath, algo)
	if err != nil {
		return err
	}

	
	if recover {
		
		showMnemonic = false
		mnemonic = ""
	}

	return printCreate(cmd, k, showMnemonic, mnemonic, outputFormat)
}

func printCreate(cmd *cobra.Command, k *keyring.Record, showMnemonic bool, mnemonic, outputFormat string) error {
	switch outputFormat {
	case keys.OutputFormatText:
		cmd.PrintErrln()
		if err := printKeyringRecord(cmd.OutOrStdout(), k, keyring.MkAccKeyOutput, outputFormat); err != nil {
			return err
		}

		
		if showMnemonic {
			if _, err := fmt.Fprintf(cmd.ErrOrStderr(),
				"\n**Important** write this mnemonic phrase in a safe place.\nIt is the only way to recover your account if you ever forget your password.\n\n%s\n\n", 
				mnemonic); err != nil {
				return fmt.Errorf("failed to print mnemonic: %v", err)
			}
		}
	case keys.OutputFormatJSON:
		out, err := keyring.MkAccKeyOutput(k)
		if err != nil {
			return err
		}

		if showMnemonic {
			out.Mnemonic = mnemonic
		}

		jsonString, err := keys.KeysCdc.MarshalJSON(out)
		if err != nil {
			return err
		}

		cmd.Println(string(jsonString))

	default:
		return fmt.Errorf("invalid output format %s", outputFormat)
	}

	return nil
}

func validateMultisigThreshold(k, nKeys int) error {
	if k <= 0 {
		return fmt.Errorf("threshold must be a positive integer")
	}
	if nKeys < k {
		return fmt.Errorf(
			"threshold k of n multisignature: %d < %d", nKeys, k)
	}
	return nil
}

type bechKeyOutFn func(k *cryptokeyring.Record) (cryptokeyring.KeyOutput, error)

func printKeyringRecord(w io.Writer, k *cryptokeyring.Record, bechKeyOut bechKeyOutFn, output string) error {
	ko, err := bechKeyOut(k)
	if err != nil {
		return err
	}

	switch output {
	case keys.OutputFormatText:
		if err := printTextRecords(w, []cryptokeyring.KeyOutput{ko}); err != nil {
			return err
		}

	case keys.OutputFormatJSON:
		out, err := keys.KeysCdc.MarshalJSON(ko)
		if err != nil {
			return err
		}

		if _, err := fmt.Fprintln(w, string(out)); err != nil {
			return err
		}
	}

	return nil
}


func printTextRecords(w io.Writer, kos []cryptokeyring.KeyOutput) error {
	out, err := yaml.Marshal(&kos)
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintln(w, string(out)); err != nil {
		return err
	}

	return nil
}