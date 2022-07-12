package caller

import (
	"math/big"
)

type DEXApiCaller interface {
	// create order
	CreateOrder(coinPair []string, amount *big.Int) (string, error)
	// get price in dex
	GetAmount(coinPair []string, coinIn *big.Int) (*big.Int, error)
	// get balance
	GetBalance(coin string) (*big.Int, error)
	// get address
	GetAddress() string
}
