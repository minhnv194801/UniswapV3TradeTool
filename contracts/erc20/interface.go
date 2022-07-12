package erc20

import (
	"math/big"
)

type Client interface {
	Approve(amount *big.Int, contractAddress string) error
	Allowance(contractAddress string) (*big.Int, error)
	GetBalance(address string) (*big.Int, error)
}
