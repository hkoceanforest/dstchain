package cli

import (
	"encoding/hex"
	"freemasonry.cc/blockchain/util"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/privval"
	"regexp"
	"strings"

	"freemasonry.cc/blockchain/x/gateway/types"
	"github.com/cosmos/cosmos-sdk/client"
)

func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Gateway transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		gatewayRegisterCmd(),
	)
	return txCmd
}

func gatewayRegisterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register [pvKey_file] [node_id] [gateway_url] [index_num]",
		Short: "register a gateway to the blockchain",
		Long: `
			register a gateway to the blockchain.
			Args:
			- pvKey_file: the file path of the private key of the validator
			- node_id: the node id of the gateway is running on
			- gateway_url: the url of the gateway
			- index_num: the index number of the gateway(max length is 7,min length is 5,only numbers are allowed)
				
			 	Example:
				stcd tx gateway register ./priv_validator_key.json bd961b42cd9286af3d74015739255a68d2d72c0e  http://127.0.0.1:50327 11111 --from=<key_or_address>
				`,
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			filePV := privval.LoadFilePVEmptyState(args[0], "")
			priv, _, err := crypto.GenerateEd25519Key(strings.NewReader(hex.EncodeToString(filePV.Key.PrivKey.Bytes())))
			if err != nil {
				return err
			}
			peerId, err := peer.IDFromPrivateKey(priv)
			wallet, err := createAccount(hex.EncodeToString(filePV.Key.PrivKey.Bytes()))
			if err != nil {
				return err
			}
			id := args[1]
			for {
				id = util.Md5String(id)
				matched, err := regexp.MatchString("^[a-z]", id)
				if err != nil {
					return err
				}
				if !matched {
					continue
				}
				end := id[:3]
				id = id + "." + end
				break
			}
			indexNum := []string{args[3]}
			msg := types.NewMsgGatewayRegister(clientCtx.GetFromAddress().String(), args[2], "", id, peerId.String(), wallet, indexNum)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
