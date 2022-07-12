package cmd

import (
	"fmt"
	"log"
	"orderbot/config"
	"orderbot/services"
	"orderbot/strategies"
	"orderbot/utils"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var (
	startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start bot",
		Run:   runStart,
	}
)

var (
	startStopBalanceStr string
	startSpendEachStr   string
	startSellCoin       string
	startBuyCoin        string
	startMinPrice       string
)

func init() {
	startCmd.Flags().StringVar(&startStopBalanceStr, "stop-balance", "", "stop balance")
	startCmd.Flags().StringVar(&startSpendEachStr, "spend-each", "", "spend each")
	startCmd.Flags().StringVar(&startSellCoin, "sell-coin", "", "sell coin name")
	startCmd.Flags().StringVar(&startBuyCoin, "buy-coin", "", "buy coin name")
	startCmd.Flags().StringVar(&startMinPrice, "min-price", "", "min price")

	startCmd.MarkFlagRequired("stop-balance")
	startCmd.MarkFlagRequired("spend-each")
	startCmd.MarkFlagRequired("sell-coin")
	startCmd.MarkFlagRequired("buy-coin")
	startCmd.MarkFlagRequired("min-price")
	RootCmd.AddCommand(startCmd)
}

func runStart(cmd *cobra.Command, args []string) {
	log.Println("------------ Start order bot ------------")
	log.Println(fmt.Sprintf("Coin sell: %s Coin buy: %s", startSellCoin, startBuyCoin))
	log.Println(fmt.Sprintf("Stop balance: %s Min price: %s", startStopBalanceStr, startMinPrice))
	log.Println(fmt.Sprintf("Spend each: %s", startSpendEachStr))
	stopBalance := utils.ParseStringTo18Decimal(startStopBalanceStr)
	spenEach := utils.ParseStringTo18Decimal(startSpendEachStr)
	minPrice := utils.ParseStringTo18Decimal(startMinPrice)
	coinPair := []string{startSellCoin, startBuyCoin}
	watcher, err := config.NewWatcherConfig()
	if err != nil {
		log.Fatal(err)
	}
	service, err := services.NewDEXOrderService(coinPair, watcher.GetDexConf(), watcher.GetCoinsMap())
	if err != nil {
		log.Fatal(err)
	}
	strategy, err := strategies.NewSplitMoneyStrategyBuilder(
		stopBalance,
		spenEach,
		minPrice,
		service,
	)
	if err != nil {
		log.Fatal(err)
	}
	if watcher.GetDexConf().TimeStep > 0 {
		err = strategy.BuildTimeStepLength(watcher.GetDexConf().TimeStep)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = strategy.Build().Start()
	if err != nil {
		log.Fatal(err)
	}
}
