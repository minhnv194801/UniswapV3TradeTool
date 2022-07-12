package caller

import (
	"math/big"
	"orderbot/config"
	"testing"
)

func TestUniswapApiCaller_CreateOrder(t *testing.T) {
	watcher, err := config.NewWatcherConfig()
	if err != nil {
		t.Error(err)
	}
	t.Log(watcher.GetDexConf())
	t.Log(watcher.GetCoinsMap())
	uniswap, err := NewUniswapApiCaller(
		watcher.GetDexConf(),
		watcher.GetCoinsMap(),
	)
	if err != nil {
		t.Error(err)
	}
	// 2000000000000000000 2DAI 50000000000000000 0.05ETH
	amount, _ := new(big.Int).SetString("2000000000000000000", 10)
	_, err = uniswap.CreateOrder([]string{"DAI", "ETH"}, amount)
	if err != nil {
		t.Error(err)
	}
}

func TestUniswapApiCaller_GetBalance(t *testing.T) {
	watcher, err := config.NewWatcherConfig()
	if err != nil {
		t.Error(err)
	}
	uniswap, err := NewUniswapApiCaller(
		watcher.GetDexConf(),
		watcher.GetCoinsMap(),
	)
	if err != nil {
		t.Error(err)
	}
	balance, err := uniswap.GetBalance("WETH")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(balance)
}

func TestUniswapApiCaller_GetAmount(t *testing.T) {
	watcher, err := config.NewWatcherConfig()
	if err != nil {
		t.Error(err)
	}
	uniswap, err := NewUniswapApiCaller(
		watcher.GetDexConf(),
		watcher.GetCoinsMap(),
	)
	if err != nil {
		t.Error(err)
	}
	amountIn, _ := new(big.Int).SetString("50000000000000000", 10)
	amount, err := uniswap.GetAmount([]string{"TA", "TB"}, amountIn)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(amount)
}
