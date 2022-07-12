package utils

import (
	"context"
	"log"
	"math/big"
	"orderbot/contracts/erc20"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

// Wrap an amount of ETH into WETH
func WrapETH(baseUrl, privateKey, wrappedAddress string, amount *big.Int) (string, error) {
	client, err := ethclient.Dial(baseUrl)
	if err != nil {
		return "", err
	}
	key, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return "", err
	}
	address := crypto.PubkeyToAddress(key.PublicKey)
	nonce, err := client.PendingNonceAt(context.Background(), address)
	if err != nil {
		return "", err
	}

	value := amount
	gasLimit := uint64(1000000) // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	toAddress := common.HexToAddress(wrappedAddress)
	transferFnSignature := []byte("deposit()")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	var data []byte
	data = append(data, methodID...)
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), key)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		panic(err)
	}
	return tx.Hash().String(), nil
}

// Convert all available WETH back to ETH
func UnWrapETH(baseUrl, privateKey, wrappedAddress string) (string, error) {
	client, err := ethclient.Dial(baseUrl)
	if err != nil {
		return "", err
	}
	key, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return "", err
	}
	address := crypto.PubkeyToAddress(key.PublicKey)
	nonce, err := client.PendingNonceAt(context.Background(), address)
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(0)
	gasLimit := uint64(1000000) // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return "", err
	}
	erc20Instance, err := erc20.NewClientErc20(baseUrl, wrappedAddress, privateKey, chainID.Int64())
	if err != nil {
		return "", err
	}
	amount, err := erc20Instance.GetBalance(address.String())
	if err != nil {
		return "", err
	}

	toAddress := common.HexToAddress(wrappedAddress)
	transferFnSignature := []byte("withdraw(uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAmount...)
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), key)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		panic(err)
	}
	return tx.Hash().String(), nil
}
