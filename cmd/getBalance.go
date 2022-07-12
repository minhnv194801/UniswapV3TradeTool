package cmd

import (
	"fmt"
	"log"
	"orderbot/config"
	"orderbot/services"
	"orderbot/utils"

	"github.com/spf13/cobra"
)

// getBalance represents the start command
var (
	getBalance = &cobra.Command{
		Use:   "balance",
		Short: "Get balance",
		Run:   runGetBalance,
	}
)

var (
	balanceCoin string
)

func init() {
	getBalance.Flags().StringVar(&balanceCoin, "coin", "", "coin name")

	getBalance.MarkFlagRequired("coin")
	RootCmd.AddCommand(getBalance)
}

func runGetBalance(cmd *cobra.Command, args []string) {
	watcher, err := config.NewWatcherConfig()
	if err != nil {
		log.Fatal(err)
	}
	service, err := services.NewDEXOrderService([]string{balanceCoin, "BNB"}, watcher.GetDexConf(), watcher.GetCoinsMap())
	if err != nil {
		log.Fatal(err)
	}

	balance, err := service.GetBalance()
	if err != nil {
		log.Fatal(err)
	}
	balanceFloat, err := utils.Parse18DecimalToFloat(balance.String())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Balance is %f \n", balanceFloat)
}
