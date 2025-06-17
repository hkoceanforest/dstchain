package cli

import (
	"github.com/spf13/cobra"

	"freemasonry.cc/blockchain/x/chat/types"
	"github.com/cosmos/cosmos-sdk/client"
)

func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "erc20 subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
	
	
	)
	return txCmd
}
