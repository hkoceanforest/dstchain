package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"freemasonry.cc/blockchain/app"
	"freemasonry.cc/blockchain/core"
	"github.com/cosmos/cosmos-sdk/client/flags"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/evmos/ethermint/encoding"
	"github.com/spf13/cobra"
	tmjson "github.com/tendermint/tendermint/libs/json"
	tmlog "github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"io"
	"os"
	"path/filepath"
)

const (
	flagTraceStore = "trace-store"
)

var encCfg = encoding.MakeConfig(app.ModuleBasics)

func ExportCmdCustom(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export-public",
		Short: "Export state to JSON",
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			homeDir, _ := cmd.Flags().GetString(flags.FlagHome)
			config.SetRoot(homeDir)

			if _, err := os.Stat(config.GenesisFile()); os.IsNotExist(err) {
				return err
			}

			db, err := openDB(config.RootDir, server.GetAppDBBackend(serverCtx.Viper))
			if err != nil {
				return err
			}

			appExporter := func(
				logger tmlog.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailAllowedAddrs []string,
				appOpts servertypes.AppOptions,
			) (servertypes.ExportedApp, error) {
				var evmosApp *app.Evmos
				homePath, ok := appOpts.Get(flags.FlagHome).(string)
				if !ok || homePath == "" {
					return servertypes.ExportedApp{}, errors.New("application home not set")
				}

				if height != -1 {
					evmosApp = app.NewEvmos(logger, db, traceStore, false, map[int64]bool{}, "", uint(1), encCfg, appOpts)

					if err := evmosApp.LoadHeight(height); err != nil {
						return servertypes.ExportedApp{}, err
					}
				} else {
					evmosApp = app.NewEvmos(logger, db, traceStore, true, map[int64]bool{}, "", uint(1), encCfg, appOpts)
				}

				return evmosApp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs)
			}

			traceWriterFile, _ := cmd.Flags().GetString(flagTraceStore)
			traceWriter, err := openTraceWriter(traceWriterFile)
			if err != nil {
				return err
			}

			height, _ := cmd.Flags().GetInt64(flags.FlagHeight)
			forZeroHeight, _ := cmd.Flags().GetBool(server.FlagForZeroHeight)
			jailAllowedAddrs, _ := cmd.Flags().GetStringSlice(server.FlagJailAllowedAddrs)

			exported, err := appExporter(serverCtx.Logger, db, traceWriter, height, forZeroHeight, jailAllowedAddrs, serverCtx.Viper)
			if err != nil {
				return fmt.Errorf("error exporting state: %v", err)
			}

			doc, err := tmtypes.GenesisDocFromFile(serverCtx.Config.GenesisFile())
			if err != nil {
				return err
			}

			doc.AppState = exported.AppState
			doc.Validators = exported.Validators
			doc.InitialHeight = exported.Height
			doc.ConsensusParams = &tmproto.ConsensusParams{
				Block: tmproto.BlockParams{
					MaxBytes:   exported.ConsensusParams.Block.MaxBytes,
					MaxGas:     exported.ConsensusParams.Block.MaxGas,
					TimeIotaMs: doc.ConsensusParams.Block.TimeIotaMs,
				},
				Evidence: tmproto.EvidenceParams{
					MaxAgeNumBlocks: exported.ConsensusParams.Evidence.MaxAgeNumBlocks,
					MaxAgeDuration:  exported.ConsensusParams.Evidence.MaxAgeDuration,
					MaxBytes:        exported.ConsensusParams.Evidence.MaxBytes,
				},
				Validator: tmproto.ValidatorParams{
					PubKeyTypes: exported.ConsensusParams.Validator.PubKeyTypes,
				},
			}

			
			appState := make(map[string]json.RawMessage)
			err = json.Unmarshal(exported.AppState, &appState)
			if err != nil {
				return err
			}

			
			if _, ok := appState["dao"]; ok {
				newDaoState, err := ExportDao(appState["dao"])
				if err != nil {
					return err
				}

				appState["dao"] = newDaoState
			}
			
			var totalStakeAmount sdk.Int
			var newStakingState json.RawMessage
			var stakeMap map[string]sdk.Int
			var vals []stakingTypes.Validator

			
			if _, ok := appState["staking"]; ok {
				newStakingState, totalStakeAmount, stakeMap, vals, err = ExportStaking(appState["staking"])
				if err != nil {
					return err
				}
				appState["staking"] = newStakingState
			}

			
			if _, ok := appState["distribution"]; ok {
				newDistributionState, err := ExportDistribution(appState["distribution"], stakeMap)
				if err != nil {
					return err
				}

				appState["distribution"] = newDistributionState
			}

			
			if _, ok := appState["slashing"]; ok {
				newSlashinState, err := ExportSlashings(appState["slashing"])
				if err != nil {
					return err
				}

				appState["slashing"] = newSlashinState

			}

			
			if _, ok := appState["bank"]; ok {
				newBankState, err := ExportBank(appState["bank"], totalStakeAmount)
				if err != nil {
					return err
				}

				appState["bank"] = newBankState
			}

			
			if _, ok := appState["gov"]; ok {
				newGovState, err := ExportGov(appState["gov"])
				if err != nil {
					return err
				}

				appState["gov"] = newGovState
			}

			
			if _, ok := appState["evm"]; ok {
				newEvmState, err := ExportEvm(appState["evm"])
				if err != nil {
					return err
				}

				appState["evm"] = newEvmState
			}

			
			if _, ok := appState["claim"]; ok {
				newClaimState, err := ExportClaim(appState["claim"])
				if err != nil {
					return err
				}
				appState["claim"] = newClaimState
			}

			newGenesisValidator := make([]tmtypes.GenesisValidator, 0)
			for _, val := range vals {
				if val.OperatorAddress == "dstvaloper1ll30h0xykgduvxxfnpy4h6yzl0770pgn7hn3lz" {
					pk, err := val.ConsPubKey()
					if err != nil {
						return err
					}
					tmPk, err := cryptocodec.ToTmPubKeyInterface(pk)
					if err != nil {
						return err
					}

					newGenesisValidator = append(newGenesisValidator, tmtypes.GenesisValidator{
						Address: sdk.ConsAddress(tmPk.Address()).Bytes(),
						PubKey:  tmPk,
						Power:   val.GetConsensusPower(sdk.DefaultPowerReduction),
						Name:    val.Description.Moniker,
					})
				}

			}

			exported.Validators = newGenesisValidator

			appStateJsonBytes, err := json.Marshal(appState)
			if err != nil {
				return err
			}

			doc.AppState = appStateJsonBytes
			doc.Validators = newGenesisValidator
			
			

			
			
			
			encoded, err := tmjson.Marshal(doc)
			if err != nil {
				return err
			}

			
			for k, v := range core.ValReplace {
				encoded = bytes.ReplaceAll(encoded, []byte(k), []byte(v))
			}

			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.OutOrStderr())
			cmd.Println(string(sdk.MustSortJSON(encoded)))
			return nil
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	cmd.Flags().Int64(server.FlagHeight, -1, "Export state from a particular height (-1 means latest height)")
	cmd.Flags().Bool(server.FlagForZeroHeight, false, "Export state to start at height zero (perform preproccessing)")
	cmd.Flags().StringSlice(server.FlagJailAllowedAddrs, []string{}, "Comma-separated list of operator addresses of jailed validators to unjail")

	return cmd
}

func openDB(rootDir string, backendType dbm.BackendType) (dbm.DB, error) {
	dataDir := filepath.Join(rootDir, "data")
	return dbm.NewDB("application", backendType, dataDir)
}

func openTraceWriter(traceWriterFile string) (w io.Writer, err error) {
	if traceWriterFile == "" {
		return
	}

	filePath := filepath.Clean(traceWriterFile)
	return os.OpenFile(
		filePath,
		os.O_WRONLY|os.O_APPEND|os.O_CREATE,
		0o600,
	)
}

type ClusterExportRecord struct {
	ClusterId    string
	DeviceAmount int
	PowerAmount  sdk.Dec
}
