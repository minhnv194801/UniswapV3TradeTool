package cmd

import (
	"fmt"
	"log"
	"orderbot/config"
	"orderbot/services"
	"orderbot/utils"

	"github.com/spf13/cobra"
)

// getPrice represents the start command
var (
	getPrice = &cobra.Command{
		Use:   "price",
		Short: "Get price",
		Run:   runGetPrice,
	}
)

var (
	priceSellCoin string
	priceBuyCoin  string
)

func init() {
	getPrice.Flags().StringVar(&priceSellCoin, "sell-coin", "", "sell coin name")
	getPrice.Flags().StringVar(&priceBuyCoin, "buy-coin", "", "buy coin name")

	getPrice.MarkFlagRequired("sell-coin")
	getPrice.MarkFlagRequired("buy-coin")
	RootCmd.AddCommand(getPrice)
}

func runGetPrice(cmd *cobra.Command, args []string) {
	coinPair := []string{priceSellCoin, priceBuyCoin}
	watcher, err := config.NewWatcherConfig()
	if err != nil {
		log.Fatal(err)
	}
	service, err := services.NewDEXOrderService(coinPair, watcher.GetDexConf(), watcher.GetCoinsMap())
	if err != nil {
		log.Fatal(err)
	}

	price, err := service.GetPrice()
	if err != nil {
		log.Fatal(err)
	}

	priceFloat, err := utils.Parse18DecimalToFloat(price.String())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Price is %f \n", priceFloat)
}
