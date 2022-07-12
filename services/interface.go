package services

import "math/big"

type Order interface {
	// create order
	CreateOrder(spend *big.Int) (string, error)
	// get price
	GetPrice() (*big.Int, error)
	// get balance
	GetBalance() (*big.Int, error)
	// get address
	GetAddress() string
	// get coin pair
	GetCoinPair() []string
}
