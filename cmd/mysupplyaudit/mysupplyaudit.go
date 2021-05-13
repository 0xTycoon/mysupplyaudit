package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/currencytycoon/mysupplyaudit"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const configDefaultPath = "./config.json"

var (
	verbose    bool
	configPath string
	conf       *mysupplyaudit.Config

	checkCmd = &cobra.Command{
		Use:   "audit [blockHeight]",
		Short: "audit the supply",
		Long:  "audit the supply of ETH by summing the Genesis block + block rewards for all succeeding blocks",
		Run:   audit,
		Args:  cobra.MinimumNArgs(0),
	}

	rootCmd = &cobra.Command{
		Use:   "mysupplyaudit",
		Short: "audit ethereum's ETH supply",
		Long:  `audit Ethereum's supply by asking your ETH full node for the blocks`,
		Run:   nil,
	}
)

func init() {
	cobra.OnInitialize()
	checkCmd.PersistentFlags().StringVarP(&configPath, "config", "c",
		configDefaultPath, "Path to the configuration file")

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false,
		"print out more debug information")
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if verbose {
			logrus.SetLevel(logrus.DebugLevel)
		} else {
			logrus.SetLevel(logrus.InfoLevel)
		}
	}
	rootCmd.AddCommand(checkCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func audit(cmd *cobra.Command, args []string) {

	s, err := mysupplyaudit.NewSupplier(configPath)
	if err != nil {
		logrus.WithError(err).Error("cannot start supplier")
		os.Exit(1)
	}

	logrus.Info("Welcome to My Supply Audit")
	var highestBlock int64
	highestBlock = -1
	if len(args) > 0 {
		highestBlock, err = strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			logrus.WithError(err).Error("invalid block number")
			os.Exit(1)
		}
		if highestBlock < 0 {
			highestBlock = 0
		}
	}
	err = s.DoAudit(highestBlock)
	if err != nil {
		logrus.WithError(err).Error("error when getting supply")
	}
}
