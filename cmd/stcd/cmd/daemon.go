package cmd

import (
	"fmt"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/core/chainnet"
	"freemasonry.cc/blockchain/core/config"
	"freemasonry.cc/log"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func daemonCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "daemon",
		Short: "secret telegram chain start",
		Run: func(cmd *cobra.Command, args []string) {
			flagLog, _ := cmd.Flags().GetBool("log")
			flagBackground, _ := cmd.Flags().GetBool("background")
			flagProxy, _ := cmd.Flags().GetString("proxy")

			programPath, _ := filepath.Abs(os.Args[0])

			runtimePath, _ := filepath.Split(programPath)

			logPath := filepath.Join(runtimePath, "log")

			if !PathExists(logPath) {
				err := os.Mkdir(logPath, 0644)
				if err != nil {
					panic(err)
				}
			}

			daemonLogPath := filepath.Join(logPath, "chain.log")

			if flagBackground {

				_, daemonStarted := os.LookupEnv("daemonStarted")

				if !daemonStarted {
					daemonArgs := os.Args[1:]
					os.Setenv("daemonStarted", "1")
					daemonCmd := exec.Command(os.Args[0], daemonArgs...)

					daemonCmd.Start()
					fmt.Println("[PID]", daemonCmd.Process.Pid)
					os.Exit(0)
				}
			}

			if flagProxy != "" {
				os.Setenv("CHAIN_PROXY", flagProxy)
				log.Info("start proxy " + flagProxy)
			}

			InitLogger(daemonLogPath, flagLog, logrus.InfoLevel)

			tendermintRepoPath := cmd.Flag("home").Value.String()

			start(tendermintRepoPath, cmd, args)
		},
	}
	cmd.Flags().Bool("log", false, "enabled log storage")
	cmd.Flags().Bool("background", false, "running in background")
	cmd.Flags().String("proxy", "", "proxy server address")
	return cmd
}

func start(cosmosRepoPath string, cmd *cobra.Command, args []string) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainClient)
	logStoraged, _ := cmd.Flags().GetBool("log")

	log.Info("log storaged:", logStoraged)
	log.Info("repo path:", cosmosRepoPath)

	var err error

	log.Info("check config dir", len(cmd.Commands()))

	if e, _ := checkSourceExist(filepath.Join(cosmosRepoPath, "config", "genesis.json")); !e {
		log.Info("init chain")

		cmd.Root().SetArgs([]string{
			"init", "node1", "--chain-id", chainnet.ChainIdDst, "--home", cosmosRepoPath,
		})
		err = cmd.Root().Execute()
		if err != nil {
			panic(err)
		}

		err = replaceConfig(cosmosRepoPath)
		if err != nil {
			panic(err)
		}
	}
	log.Info("check genesis.json")

	err = checkGenesisFile(cosmosRepoPath)
	if err != nil {
		panic(err)
	}

	log.Info("check config.toml")

	err = checkConfigFile(cosmosRepoPath)
	if err != nil {
		panic(err)
	}

	err = checkAppToml(cosmosRepoPath)
	if err != nil {
		panic(err)
	}

	err = checkClientToml(cosmosRepoPath)
	if err != nil {
		panic(err)
	}

	err = checkDataDir(cosmosRepoPath)
	if err != nil {
		panic(err)
	}

	err = checkValidatorStateJson(cosmosRepoPath)
	if err != nil {
		panic(err)
	}

	log.WithField("path", cosmosRepoPath).Info("chain repo")

	cmd.Root().SetArgs([]string{
		"start", "--log_format", "json", "--home", cosmosRepoPath,
	})

	logLevelSet, ok := os.LookupEnv("DST_LOGGING")
	if ok {
		_, err := logrus.ParseLevel(logLevelSet)
		if err == nil {

			cmd.Root().SetArgs([]string{
				"start", "--log_format", "json", "--log_level", logLevelSet, "--home", cosmosRepoPath,
			})
		}
	}

	log.Info("start chain")

	err = cmd.Root().Execute()
	if err != nil {
		log.WithError(err).Error("start.RunE")
		panic(err)
	}

	log.Info("exit")
}

func checkDataDir(path string) error {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainClient)
	dataPath := filepath.Join(path, "data")

	if exist, _ := checkSourceExist(dataPath); exist {
		return nil
	}
	err := os.Mkdir(dataPath, os.ModePerm)
	if err != nil {
		log.WithError(err).WithField("jpath", dataPath).Error("o.Mkdir")
	}
	return err
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}

func checkGenesisFile(path string) error {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainClient)
	genesisPath := filepath.Join(path, "config", "genesis.json")

	if exist, _ := checkSourceExist(genesisPath); exist {
		return nil
	}
	err := ioutil.WriteFile(genesisPath, []byte(config.GenesisJson), os.ModePerm)
	if err != nil {
		log.WithError(err).WithField("file", genesisPath).Error("ioutil.WriteFile")
	}
	return err
}

func checkConfigFile(path string) error {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainClient)
	configPath := filepath.Join(path, "config", "config.toml")

	if exist, _ := checkSourceExist(configPath); exist {
		return nil
	}
	err := ioutil.WriteFile(configPath, []byte(config.ConfigToml), os.ModePerm)
	if err != nil {
		log.WithError(err).WithField("file", configPath).Error("ioutil.WriteFile")
	}
	return err
}

func checkAppToml(path string) error {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainClient)
	filePath := filepath.Join(path, "config", "app.toml")

	if exist, _ := checkSourceExist(filePath); exist {
		return nil
	}
	err := ioutil.WriteFile(filePath, []byte(config.AppToml), os.ModePerm)
	if err != nil {
		log.WithError(err).WithField("file", filePath).Error("ioutil.WriteFile")
	}
	return err
}

func checkClientToml(path string) error {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainClient)
	filePath := filepath.Join(path, "config", "client.toml")

	if exist, _ := checkSourceExist(filePath); exist {
		return nil
	}
	err := ioutil.WriteFile(filePath, []byte(config.ClientToml), os.ModePerm)
	if err != nil {
		log.WithError(err).WithField("file", filePath).Error("ioutil.WriteFile")
	}
	return err
}

func checkValidatorStateJson(path string) error {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainClient)
	filePath := filepath.Join(path, "data", "priv_validator_state.json")

	if exist, _ := checkSourceExist(filePath); exist {
		return nil
	}
	err := ioutil.WriteFile(filePath, []byte(config.ValidatorStateJson), os.ModePerm)
	if err != nil {
		log.WithError(err).WithField("file", filePath).Error("ioutil.WriteFile")
	}
	return err
}

func replaceConfig(path string) error {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainClient)
	appTomlPath := filepath.Join(path, "config", "app.toml")
	if exist, err := checkSourceExist(appTomlPath); !exist {
		return err
	}
	genesisPath := filepath.Join(path, "config", "genesis.json")
	if exist, err := checkSourceExist(genesisPath); !exist {
		return err
	}
	configTomlPath := filepath.Join(path, "config", "config.toml")
	if exist, err := checkSourceExist(configTomlPath); !exist {
		return err
	}
	clientTomlPath := filepath.Join(path, "config", "client.toml")
	if exist, err := checkSourceExist(clientTomlPath); !exist {
		return err
	}
	appTomlBuf, err := ioutil.ReadFile(appTomlPath)
	if err != nil {
		log.WithError(err).WithField("path", appTomlPath).Error("ioutil.ReadFile")
		return err
	}
	genesisBuf, err := ioutil.ReadFile(genesisPath)
	if err != nil {
		log.WithError(err).WithField("path", genesisPath).Error("ioutil.ReadFile")
		return err
	}
	configTomlBuf, err := ioutil.ReadFile(configTomlPath)
	if err != nil {
		log.WithError(err).WithField("path", configTomlPath).Error("ioutil.ReadFile")
		return err
	}
	clientTomlBuf, err := ioutil.ReadFile(clientTomlPath)
	if err != nil {
		log.WithError(err).WithField("path", clientTomlPath).Error("ioutil.ReadFile")
		return err
	}

	appTomlContent := string(appTomlBuf)
	genesisContent := string(genesisBuf)
	configTomlContent := string(configTomlBuf)
	clientTomlContent := string(clientTomlBuf)

	if appTomlContent != config.AppToml {
		err = ioutil.WriteFile(appTomlPath, []byte(config.AppToml), os.ModePerm)
		if err != nil {
			log.WithError(err).WithField("path", appTomlPath).Error("ioutil.WriteFile")
			return err
		}
	}
	if genesisContent != config.GenesisJson {
		err = ioutil.WriteFile(genesisPath, []byte(config.GenesisJson), os.ModePerm)
		if err != nil {
			log.WithError(err).WithField("path", genesisPath).Error("ioutil.WriteFile")
			return err
		}
	}

	if configTomlContent != config.ConfigToml {
		err = ioutil.WriteFile(configTomlPath, []byte(config.ConfigToml), os.ModePerm)
		if err != nil {
			log.WithError(err).WithField("path", configTomlPath).Error("ioutil.WriteFile")
			return err
		}
	}

	if clientTomlContent != config.ClientToml {
		err = ioutil.WriteFile(clientTomlPath, []byte(config.ClientToml), os.ModePerm)
		if err != nil {
			log.WithError(err).WithField("path", clientTomlPath).Error("ioutil.WriteFile")
			return err
		}
	}
	return err
}

func checkSourceExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func InitLogger(logPath string, enableLogSave bool, defaultLevel logrus.Level) {

	log.InitLogger(defaultLevel)

	if enableLogSave {
		logSaveTime := time.Hour * 24 * 7
		logSplitTime := time.Hour * 24
		log.EnableLogStorage(logPath, logSaveTime, logSplitTime)
	}

	logDefault, ok := os.LookupEnv("DST_LOGGING")
	if ok {
		logDefault = strings.ReplaceAll(logDefault, "\"", "")
		level, err := logrus.ParseLevel(logDefault)
		fmt.Println("env define default-log-level:", logDefault)
		if err == nil {
			log.ResetAllModuleLevel(level)
			fmt.Println("used default-log-level:", level)
		}
	}

	chainLogging, ok := os.LookupEnv("DST_LOGGING")
	if ok {
		chainLogging = strings.Trim(chainLogging, `"`)
		sets := strings.Split(chainLogging, ";")
		for _, set := range sets {

			kvs := strings.Split(set, ":")

			if len(kvs) != 2 {
				continue
			}
			moduleStr := kvs[0]
			levelStr := kvs[1]

			logModule := log.LogModule(moduleStr)

			if !log.CheckModuleExists(logModule) {
				continue
			}
			level, err := logrus.ParseLevel(levelStr)
			if err == nil {
				log.SetModuleLevel(logModule, level)
			}
		}
		fmt.Println("module-log-level--->")
		var keys []string

		models := log.GetModelLevels()
		for k, _ := range models {
			keys = append(keys, string(k))
		}
		sort.Strings(keys)

		for _, moduleName := range keys {
			fmt.Println(moduleName, ":", models[log.LogModule(moduleName)])
		}
		fmt.Println("<---")
	}
}

type CmdOutput struct {
	logWriter *rotatelogs.RotateLogs
	lastMsg   []byte
}

func (this *CmdOutput) Write(p []byte) (n int, err error) {
	this.lastMsg = p
	return this.logWriter.Write(p)
}

func (this *CmdOutput) GetLastMsg() []byte {
	return this.lastMsg
}

func NewCmdOutput(baseLogPath string, life, split time.Duration) (*CmdOutput, error) {
	writer, err := rotatelogs.New(
		baseLogPath+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(baseLogPath),
		rotatelogs.WithMaxAge(life),
		rotatelogs.WithRotationTime(split),
	)
	if err != nil {
		return nil, err
	}
	cmdOutput := CmdOutput{
		logWriter: writer,
	}
	return &cmdOutput, nil
}
