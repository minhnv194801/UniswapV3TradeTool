package erc20

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ClientErc20 struct {
	caller  *Erc20
	privK   *ecdsa.PrivateKey
	chainId int64
}

func (c *ClientErc20) Approve(amount *big.Int, contractAddress string) error {
	opt, err := bind.NewKeyedTransactorWithChainID(c.privK, big.NewInt(c.chainId))
	if err != nil {
		return err
	}
	opt.GasLimit = 1000000
	_, err = c.caller.Approve(opt, common.HexToAddress(contractAddress), amount)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientErc20) Allowance(contractAddress string) (*big.Int, error) {
	opt, err := bind.NewKeyedTransactorWithChainID(c.privK, big.NewInt(c.chainId))
	if err != nil {
		return nil, err
	}
	opt.GasLimit = 1000000
	allowance, err := c.caller.Allowance(nil, opt.From, common.HexToAddress(contractAddress))
	if err != nil {
		return nil, err
	}
	return allowance, nil
}

func (c *ClientErc20) GetBalance(address string) (*big.Int, error) {
	balance, err := c.caller.BalanceOf(nil, common.HexToAddress(address))
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func NewClientErc20(baseUrl, address, privK string, chainId int64) (Client, error) {
	client, err := ethclient.Dial(baseUrl)
	if err != nil {
		return nil, err
	}
	caller, err := NewErc20(common.HexToAddress(address), client)
	if err != nil {
		return nil, err
	}
	key, err := crypto.HexToECDSA(privK)
	if err != nil {
		return nil, err
	}
	return &ClientErc20{
		caller:  caller,
		privK:   key,
		chainId: chainId,
	}, nil
}
