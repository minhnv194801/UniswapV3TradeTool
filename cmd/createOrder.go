package cmd

import (
	"fmt"
	"log"
	"orderbot/config"
	"orderbot/services"
	"orderbot/utils"

	"github.com/spf13/cobra"
)

// createOrder represents the start command
var (
	createOrder = &cobra.Command{
		Use:   "order",
		Short: "Create order",
		Run:   runCreateOrder,
	}
)

var (
	orderSellCoin  string
	orderBuyCoin   string
	orderSpendCoin string
)

func init() {
	createOrder.Flags().StringVar(&orderSellCoin, "sell-coin", "", "sell coin name")
	createOrder.Flags().StringVar(&orderBuyCoin, "buy-coin", "", "buy coin name")
	createOrder.Flags().StringVar(&orderSpendCoin, "spend", "", "spend coin")

	createOrder.MarkFlagRequired("sell-coin")
	createOrder.MarkFlagRequired("buy-coin")
	createOrder.MarkFlagRequired("spend")

	RootCmd.AddCommand(createOrder)
}

func runCreateOrder(cmd *cobra.Command, args []string) {
	coinPair := []string{orderSellCoin, orderBuyCoin}
	watcher, err := config.NewWatcherConfig()
	if err != nil {
		log.Fatal(err)
	}
	service, err := services.NewDEXOrderService(coinPair, watcher.GetDexConf(), watcher.GetCoinsMap())
	if err != nil {
		log.Fatal(err)
	}
	spendable := utils.ParseStringTo18Decimal(orderSpendCoin)
	fmt.Printf("Spend %s to buy coin\n", orderSpendCoin)
	txHash, err := service.CreateOrder(spendable)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Success create order, txHash=" + string(txHash) + " \n")
}
