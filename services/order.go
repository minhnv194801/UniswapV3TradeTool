package services

import (
	"math/big"
	"orderbot/caller"
	"orderbot/models"
)

type OrderService struct {
	coinPair  []string
	apiCaller caller.DEXApiCaller
}

func (o *OrderService) GetAddress() string {
	return o.apiCaller.GetAddress()
}

func (o *OrderService) CreateOrder(spend *big.Int) (string, error) {
	return o.apiCaller.CreateOrder(o.coinPair, spend)
}

func (o *OrderService) GetPrice() (*big.Int, error) {
	return o.apiCaller.GetAmount(o.coinPair, big.NewInt(1000000000000000000))
}

func (o *OrderService) GetCoinPair() []string {
	return o.coinPair
}

func (o *OrderService) GetBalance() (*big.Int, error) {
	return o.apiCaller.GetBalance(o.coinPair[0])
}

func NewDEXOrderService(coinPair []string, setting models.DEXSetting, coinMap models.CoinsMap) (Order, error) {
	client, err := caller.NewUniswapApiCaller(setting, coinMap)
	if err != nil {
		return nil, err
	}
	return &OrderService{
		coinPair:  coinPair,
		apiCaller: client,
	}, nil
}
