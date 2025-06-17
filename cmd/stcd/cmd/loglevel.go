package cmd

import (
	"fmt"
	"freemasonry.cc/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"sort"
)

func logLevelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "log-level",
		Short: "log module level info",
		Long:  `View the log level details of the current blockchain`,
		Example: "The steps to modify the blockchain log level are as follows:" +
			"\n1.Set environment variables `export DST_LOGGING=\"defaultLogLevel;ModuleName1:logLevel;ModuleName2:logLevel\"` (Module log level is optional)" +
			"\n2.Run the command to observe if it can take effect `stcd log-level`" +
			"\n3.Restart stcd",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("*********************************************")
			fmt.Println("*************   log level info  *************")
			fmt.Println("*********************************************")
			InitLogger("", false, logrus.InfoLevel)
			_, ok := os.LookupEnv("DST_LOGGING")
			if !ok {
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

		},
	}
	return cmd
}
