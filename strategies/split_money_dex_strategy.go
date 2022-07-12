package strategies

import (
	"errors"
	"fmt"
	"math/big"
	"orderbot/consts"
	"orderbot/filewriter"
	"orderbot/logger"
	"orderbot/services"
	"orderbot/utils"
	"sync"
	"time"
)

type SplitMoneyStrategy struct {
	// time between each order, exp: 20
	timeStepLength int64
	// coin remain that cancel process immediately, exp: 1 * 10^18
	stopBalance *big.Int
	// coin use each order exp: 0.1 * 10^18
	spendEach *big.Int
	// order service
	orderService services.Order
	// logger
	logger logger.Logger
	// isRunning
	isRunning bool
	// minPrice
	minPrice *big.Int
	// sync
	syncMutex sync.RWMutex
	// writer
	fileWrite filewriter.WriterWriteLine
}

func (strategy *SplitMoneyStrategy) Start() error {
	// start epoch
	c := time.NewTicker(time.Duration(strategy.timeStepLength) * time.Second)
	defer c.Stop()
	errChan := make(chan error)
	for {
		select {
		case <-c.C:
			if strategy.isRunning {
				break
			}
			go func() {
				strategy.SetRunning(true)
				defer strategy.SetRunning(false)
				remain, err := strategy.orderService.GetBalance()
				if err != nil {
					strategy.logger.Warn("SplitMoneyStrategy.Start", "Failed to get balance", 1)
					remain, err = strategy.orderService.GetBalance()
					if err != nil {
						strategy.logger.Error("SplitMoneyStrategy.Start", "Failed to get balance", 2)
						errChan <- err
						return
					}
				}
				remainFloat, err := utils.Parse18DecimalToFloat(remain.String())
				if err != nil {
					errChan <- err
					return
				}
				strategy.logger.Info("SplitMoneyStrategy.Start", "Coin remain is: "+fmt.Sprint(remainFloat), 1)
				// if coin is low, stop
				if remain.Cmp(strategy.stopBalance) < 0 {
					errChan <- errors.New(consts.ErrOutOfCoins)
					return
				}
				price, err := strategy.orderService.GetPrice()
				if err != nil {
					strategy.logger.Warn("SplitMoneyStrategy.Start", "Failed to get price", 1)
					price, err = strategy.orderService.GetPrice()
					if err != nil {
						strategy.logger.Error("SplitMoneyStrategy.Start", "Failed to get price", 2)
						errChan <- err
						return
					}
				}
				priceFloat, err := utils.Parse18DecimalToFloat(price.String())
				if err != nil {
					errChan <- err
					return
				}
				strategy.logger.Info("SplitMoneyStrategy.Start", fmt.Sprintf("Price is: %f", priceFloat), 1)
				// compare to min price
				if price.Cmp(strategy.minPrice) < 0 {
					// lower than min price
					strategy.logger.Info("SplitMoneyStrategy.Start", fmt.Sprintf("Price is too low, abort: %f", priceFloat), 1)
					errChan <- errors.New("price is too low")
					return
				}
				txHash, err := strategy.orderService.CreateOrder(strategy.spendEach)
				if err != nil {
					strategy.logger.Warn("SplitMoneyStrategy.Start", "Failed to place order", 1)
					txHash, err = strategy.orderService.CreateOrder(strategy.spendEach)
					if err != nil {
						strategy.logger.Error("SplitMoneyStrategy.Start", "Failed to place order", 2)
						errChan <- err
						return
					}
				}
				spendFloat, err := utils.Parse18DecimalToFloat(strategy.spendEach.String())
				if err != nil {
					errChan <- err
					return
				}
				err = strategy.fileWrite.WriteData([]string{txHash,
					strategy.orderService.GetCoinPair()[0],
					strategy.orderService.GetCoinPair()[1],
					fmt.Sprint(spendFloat),
				})
				if err != nil {
					strategy.logger.Warn("SplitMoneyStrategy.Start", "Failed to write file", 1)
				}
				strategy.logger.Info("SplitMoneyStrategy.Start",
					fmt.Sprintf("Success create order, spend=%f txhash=%s", spendFloat, txHash), 1)
			}()
		case err := <-errChan:
			return err
		}
	}
}

type SplitMoneyStrategyBuilder struct {
	strategy *SplitMoneyStrategy
}

func (strategy *SplitMoneyStrategy) SetRunning(isRunning bool) {
	strategy.syncMutex.Lock()
	defer strategy.syncMutex.Unlock()
	strategy.isRunning = isRunning
}

func NewSplitMoneyStrategyBuilder(stopBalance *big.Int, spendEach, minPrice *big.Int, service services.Order) (*SplitMoneyStrategyBuilder, error) {
	zapLogger, err := logger.NewZapLogger(consts.LogInitial, consts.LogThereafter, "OrderBot/SplitMoney")
	if err != nil {
		return nil, err
	}
	return &SplitMoneyStrategyBuilder{
		strategy: &SplitMoneyStrategy{
			timeStepLength: consts.SplitMoneyDefaultTimeLength,
			stopBalance:    stopBalance,
			spendEach:      spendEach,
			minPrice:       minPrice,
			logger:         zapLogger,
			orderService:   service,
			isRunning:      false,
			syncMutex:      sync.RWMutex{},
			fileWrite:      filewriter.NewCSVWriter("./data/", service.GetAddress(), []string{"TxnHash", "Coin-Sell", "Coin-Buy", "Spend"}),
		},
	}, nil
}

func (builder *SplitMoneyStrategyBuilder) BuildTimeStepLength(stepLength int64) error {
	if stepLength <= 0 {
		return errors.New(consts.ErrBadRequest)
	}
	builder.strategy.timeStepLength = stepLength

	return nil
}

func (builder *SplitMoneyStrategyBuilder) BuildStopBalance(stopBalance *big.Int) error {
	if stopBalance.Cmp(big.NewInt(0)) <= 0 {
		return errors.New(consts.ErrBadRequest)
	}
	builder.strategy.stopBalance = stopBalance

	return nil
}

func (builder *SplitMoneyStrategyBuilder) BuildSpendEach(spendEach *big.Int) error {
	if spendEach.Cmp(big.NewInt(0)) <= 0 {
		return errors.New(consts.ErrBadRequest)
	}
	builder.strategy.spendEach = spendEach

	return nil
}

func (builder *SplitMoneyStrategyBuilder) BuildMinPrice(minPrice *big.Int) error {
	if minPrice.Cmp(big.NewInt(0)) <= 0 {
		return errors.New(consts.ErrBadRequest)
	}
	builder.strategy.spendEach = minPrice

	return nil
}

func (builder *SplitMoneyStrategyBuilder) Build() Strategy {
	return builder.strategy
}
