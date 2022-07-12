package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var RootCmd = &cobra.Command{
	Use:   "order-bot",
	Short: "Order bot is a tool to get price, get balance, or start automatic place order bot",
	Long: `Order bot is a tool to get price, get balance, or start automatic place order bot`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

}


