package caller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"orderbot/consts"
	"orderbot/contracts/erc20"
	"orderbot/models"
	"orderbot/utils"
	"time"

	sdkCoreEntities "github.com/daoleno/uniswap-sdk-core/entities"
	sdkConstant "github.com/daoleno/uniswapv3-sdk/constants"
	"github.com/daoleno/uniswapv3-sdk/entities"
	"github.com/daoleno/uniswapv3-sdk/examples/contract"
	"github.com/daoleno/uniswapv3-sdk/periphery"
	sdkUtils "github.com/daoleno/uniswapv3-sdk/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type UniswapApiCaller struct {
	contractAddress string
	quoterAddress   string
	factoryAddress  string
	slippage        float64
	coinbase        string
	wrapToken       string
	rpcUrl          string
	privK           string
	chainId         int64
	coinMap         map[string]string
}

func (u *UniswapApiCaller) GetAddress() string {
	key, err := crypto.HexToECDSA(u.privK)
	if err != nil {
		return ""
	}
	address := crypto.PubkeyToAddress(key.PublicKey)
	return address.String()
}

func (u *UniswapApiCaller) CreateOrder(coinPair []string, amount *big.Int) (string, error) {
	// check if swap coinbase
	if len(coinPair) != 2 {
		return "", errors.New(consts.ErrBadRequest)
	}
	// get contract
	sellContract := ""
	buyContract := ""
	txHash := ""
	if coinPair[0] == u.coinbase {
		ok := false
		sellContract, ok = u.coinMap[u.wrapToken]
		if !ok {
			return "", errors.New(consts.ErrSettingInvalid + " " + u.wrapToken)
		}
	} else {
		ok := false
		sellContract, ok = u.coinMap[coinPair[0]]
		if !ok {
			return "", errors.New(consts.ErrSettingInvalid + " " + coinPair[0])
		}
	}
	if coinPair[1] == u.coinbase {
		ok := false
		buyContract, ok = u.coinMap[u.wrapToken]
		if !ok {
			return "", errors.New(consts.ErrSettingInvalid + " " + u.wrapToken)
		}
	} else {
		ok := false
		buyContract, ok = u.coinMap[coinPair[1]]
		if !ok {
			return "", errors.New(consts.ErrSettingInvalid + " " + coinPair[1])
		}
	}

	if coinPair[0] == u.coinbase {
		// if coinbase to token
		// wrap coin
		fmt.Println("Wrapping coinbase...")
		utils.WrapETH(u.rpcUrl, u.privK, sellContract, amount)
		time.Sleep(500 * time.Millisecond)
	}

	// approve
	fmt.Println("Approving sell coins...")
	erc20Client, err := erc20.NewClientErc20(u.rpcUrl, sellContract, u.privK, u.chainId)
	if err != nil {
		return "", err
	}
	err = erc20Client.Approve(amount, u.contractAddress)
	if err != nil {
		return "", err
	}
	time.Sleep(500 * time.Millisecond)
	allowance, err := erc20Client.Allowance(u.contractAddress)
	if err != nil {
		return "", err
	}
	attempt := 0
	for allowance.Cmp(amount) < 0 && attempt < consts.MaxAttempt {
		attempt++
		time.Sleep(500 * time.Millisecond)
		allowance, err = erc20Client.Allowance(u.contractAddress)
		if err != nil {
			return "", err
		}
	}
	if attempt == consts.MaxAttempt {
		return "", errors.New(consts.ErrMaxAttemptReach)
	}
	// done approve
	// swap
	fmt.Println("Performing swap...")
	txHash, err = u.swap(sellContract, buyContract, amount)
	if err != nil {
		return "", err
	}

	//Unwrap all owned wrapped eth in case it is token to coinbase
	if coinPair[1] == u.coinbase {
		fmt.Println("Unwrapping all owned WETH...")
		erc20Client, err := erc20.NewClientErc20(u.rpcUrl, buyContract, u.privK, u.chainId)
		if err != nil {
			return "", err
		}
		balance, err := erc20Client.GetBalance(u.GetAddress())
		if err != nil {
			return "", err
		}
		attempt := 0
		// Run a loop to wait for WETH arrived and then unwrap it
		for balance.Cmp(big.NewInt(0)) == 0 && attempt < consts.MaxAttempt {
			attempt++
			time.Sleep(500 * time.Millisecond)
			balance, err = erc20Client.GetBalance(u.GetAddress())
			if err != nil {
				return "", err
			}
		}
		if attempt == consts.MaxAttempt {
			return "", errors.New(consts.ErrMaxAttemptReach)
		}
		utils.UnWrapETH(u.rpcUrl, u.privK, buyContract)
	}
	return txHash, nil
}

func (u *UniswapApiCaller) GetAmount(coinPair []string, coinIn *big.Int) (*big.Int, error) {
	if len(coinPair) != 2 {
		return nil, errors.New(consts.ErrBadRequest)
	}
	sellAddress := ""
	buyAddress := ""
	if coinPair[0] == u.coinbase {
		sellAddress = u.coinMap[u.wrapToken]
	} else {
		ok := false
		if sellAddress, ok = u.coinMap[coinPair[0]]; !ok {
			return nil, errors.New(consts.ErrNoCoinFound)
		}
	}
	if coinPair[1] == u.coinbase {
		buyAddress = u.coinMap[u.wrapToken]
	} else {
		ok := false
		if buyAddress, ok = u.coinMap[coinPair[1]]; !ok {
			return nil, errors.New(consts.ErrNoCoinFound)
		}
	}

	tokenInAddress := common.HexToAddress(sellAddress)
	tokenOutAddress := common.HexToAddress(buyAddress)
	client, err := ethclient.Dial(u.rpcUrl)
	if err != nil {
		panic(err)
	}
	quoterAddress := common.HexToAddress(u.quoterAddress)
	quoterContract, err := contract.NewUniswapv3Quoter(quoterAddress, client)
	if err != nil {
		panic(err)
	}

	sqrtPriceLimitX96 := big.NewInt(0)
	var bestAmountOut = big.NewInt(0)
	poolFees := []int{int(sdkConstant.FeeLow), int(sdkConstant.FeeMedium), int(sdkConstant.FeeHigh)}
	for _, poolFee := range poolFees {
		var out []interface{}
		rawCaller := &contract.Uniswapv3QuoterRaw{Contract: quoterContract}
		err = rawCaller.Call(nil, &out, "quoteExactInputSingle", tokenInAddress, tokenOutAddress,
			big.NewInt(int64(poolFee)), coinIn, sqrtPriceLimitX96)
		if err == nil {
			amountOut := out[0].(*big.Int)
			if bestAmountOut.Cmp(amountOut) == -1 {
				bestAmountOut = amountOut
			}
		}
	}
	if bestAmountOut.Cmp(big.NewInt(0)) == 0 {
		return bestAmountOut, errors.New("no pool found for the token pair")
	}
	return bestAmountOut, nil
}

func (u *UniswapApiCaller) GetBalance(coin string) (*big.Int, error) {
	key, err := crypto.HexToECDSA(u.privK)
	if err != nil {
		return nil, err
	}
	address := crypto.PubkeyToAddress(key.PublicKey)

	if coin == u.coinbase {
		// if is coinbase
		client := http.DefaultClient
		data := map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  "eth_getBalance",
			"params": []string{
				address.String(),
				"latest",
			},
			"id": 1,
		}
		b, _ := json.Marshal(data)
		rq, err := http.NewRequest(http.MethodPost, u.rpcUrl, bytes.NewBuffer(b))
		if err != nil {
			return nil, err
		}
		rq.Header.Set("Content-Type", "application/json")
		res, err := client.Do(rq)
		if err != nil {
			return nil, err
		}
		d, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		ret := map[string]interface{}{}
		err = json.Unmarshal(d, &ret)
		if err != nil {
			return nil, err
		}
		retInt, _ := new(big.Int).SetString(ret["result"].(string)[2:], 16)
		return retInt, nil
	}
	if _, ok := u.coinMap[coin]; !ok {
		return nil, errors.New(consts.ErrBadRequest)
	}
	erc20Client, err := erc20.NewClientErc20(u.rpcUrl, u.coinMap[coin], u.privK, u.chainId)
	if err != nil {
		return nil, err
	}
	balance, err := erc20Client.GetBalance(address.String())
	if err != nil {
		return nil, err
	}
	return balance, nil
}

// Perform swaping between two tokens
func (u *UniswapApiCaller) swap(sellContract, buyContract string, amount *big.Int) (string, error) {
	client, err := ethclient.Dial(u.rpcUrl)
	if err != nil {
		return "", err
	}
	tokenInAddress := common.HexToAddress(sellContract)
	tokenOutAddress := common.HexToAddress(buyContract)
	chainId, err := client.NetworkID(context.Background())
	if err != nil {
		return "", err
	}
	//Ignore decimal, name, symbol
	tokenIn := sdkCoreEntities.NewToken(uint(chainId.Uint64()), tokenInAddress, 18, "TA", "Token A")
	tokenOut := sdkCoreEntities.NewToken(uint(chainId.Uint64()), tokenOutAddress, 18, "TB", "Token B")
	pool, err := u.getBestPool(tokenIn, tokenOut)
	if err != nil {
		return "", err
	}
	numerator := u.slippage
	denominator := big.NewInt(1)
	for numerator < 1 {
		denominator.Mul(denominator, big.NewInt(10))
		numerator *= 10
	}
	slippageTolerance := sdkCoreEntities.NewPercent(big.NewInt(int64(numerator)), denominator)
	//after 5 minutes
	d := time.Now().Add(time.Minute * time.Duration(15)).Unix()
	deadline := big.NewInt(d)

	r, err := entities.NewRoute([]*entities.Pool{pool}, tokenIn, tokenOut)
	if err != nil {
		return "", err
	}

	trade, err := entities.FromRoute(r, sdkCoreEntities.FromRawAmount(tokenIn, amount), sdkCoreEntities.ExactInput)
	if err != nil {
		return "", err
	}

	fromAddress := common.HexToAddress(u.GetAddress())
	params, err := periphery.SwapCallParameters([]*entities.Trade{trade}, &periphery.SwapOptions{
		SlippageTolerance: slippageTolerance,
		Recipient:         fromAddress,
		Deadline:          deadline,
	})
	if err != nil {
		return "", err
	}
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", err
	}
	value := big.NewInt(0)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}
	swapRouterAddress := common.HexToAddress(u.contractAddress)
	gasLimit := uint64(1000000)
	tx := types.NewTransaction(nonce, swapRouterAddress, value, gasLimit, gasPrice, params.Calldata)
	key, err := crypto.HexToECDSA(u.privK)
	if err != nil {
		return "", err
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), key)
	if err != nil {
		return "", err
	}
	txHash := signedTx.Hash().String()
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}
	return txHash, nil
}

// Get the address for the pool of a token pair with the correspond fee
func (u *UniswapApiCaller) getPoolAddress(token0, token1 common.Address, fee *big.Int) (common.Address, error) {
	client, err := ethclient.Dial(u.rpcUrl)
	if err != nil {
		return common.Address{}, err
	}
	factoryAddress := common.HexToAddress(u.factoryAddress)
	factoryContract, err := contract.NewUniswapv3Factory(factoryAddress, client)
	if err != nil {
		return common.Address{}, err
	}
	poolAddress, err := factoryContract.GetPool(nil, token0, token1, fee)
	if err != nil {
		return common.Address{}, err
	}
	if poolAddress == (common.Address{}) {
		return common.Address{}, err
	}
	return poolAddress, nil
}

// Create pool entity of a token pair with the correspond fee
func (u *UniswapApiCaller) getPool(token0, token1 *sdkCoreEntities.Token, fee *big.Int) (*entities.Pool, error) {
	client, err := ethclient.Dial(u.rpcUrl)
	if err != nil {
		return nil, err
	}
	poolAddress, err := u.getPoolAddress(token0.Address, token1.Address, fee)
	if err != nil {
		return nil, err
	}
	poolContract, err := contract.NewUniswapv3Pool(poolAddress, client)
	if err != nil {
		return nil, err
	}
	liquidity, err := poolContract.Liquidity(nil)
	if err != nil {
		return nil, err
	}
	slot0, err := poolContract.Slot0(nil)
	if err != nil {
		return nil, err
	}
	pooltickMin, err := poolContract.Ticks(nil, big.NewInt(sdkUtils.MinTick))
	if err != nil {
		return nil, err
	}
	pooltickMax, err := poolContract.Ticks(nil, big.NewInt(sdkUtils.MaxTick))
	if err != nil {
		return nil, err
	}
	feeAmount := sdkConstant.FeeAmount(fee.Uint64())
	tickSpacing, err := poolContract.TickSpacing(nil)
	if err != nil {
		return nil, err
	}
	ticks := []entities.Tick{
		{
			Index: entities.NearestUsableTick(sdkUtils.MinTick,
				int(tickSpacing.Int64())),
			LiquidityNet:   pooltickMin.LiquidityNet,
			LiquidityGross: pooltickMin.LiquidityGross,
		},
		{
			Index: entities.NearestUsableTick(sdkUtils.MaxTick,
				int(tickSpacing.Int64())),
			LiquidityNet:   pooltickMax.LiquidityNet,
			LiquidityGross: pooltickMax.LiquidityGross,
		},
	}
	p, err := entities.NewTickListDataProvider(ticks, sdkConstant.TickSpacings[feeAmount])
	if err != nil {
		return nil, err
	}
	return entities.NewPool(token0, token1, feeAmount, slot0.SqrtPriceX96, liquidity, int(slot0.Tick.Int64()), p)
}

// Create the best pool for a token pair
//
// The function will check for pool with 0.05, 0.3 and 1% fee and return the pool with the best amount out
func (u *UniswapApiCaller) getBestPool(token0, token1 *sdkCoreEntities.Token) (*entities.Pool, error) {
	var bestPool *entities.Pool
	poolFees := []int{int(sdkConstant.FeeLow), int(sdkConstant.FeeMedium), int(sdkConstant.FeeHigh)}
	for _, poolFee := range poolFees {
		pool, err := u.getPool(token0, token1, big.NewInt(int64(poolFee)))
		if err == nil {
			if bestPool == nil {
				bestPool = pool
			} else {
				if pool.Token0Price().GreaterThan(bestPool.Token0Price().Fraction) {
					bestPool = pool
				}
			}
		}
	}
	if bestPool == nil {
		return nil, errors.New("no pool available for token pair")
	}
	return bestPool, nil
}

func NewUniswapApiCaller(
	setting models.DEXSetting,
	mapCoins map[string]string,
) (DEXApiCaller, error) {
	return &UniswapApiCaller{
		coinbase:        setting.CoinbaseName,
		wrapToken:       setting.WrapToken,
		slippage:        setting.Slippage,
		rpcUrl:          setting.BaseUrl,
		privK:           setting.PrivateKey,
		chainId:         setting.ChainId,
		coinMap:         mapCoins,
		contractAddress: setting.ContractAddress,
		quoterAddress:   setting.QuoterAddress,
		factoryAddress:  setting.FactoryAddress,
	}, nil
}

func UniswapApiCallerTest(
	contractAddress string,
	slippage float64,
	coinbase string,
	wrapToken string,
	rpcUrl string,
	privK string,
	chainId int64,
	coinMap map[string]string,
) (DEXApiCaller, error) {
	return &UniswapApiCaller{
		coinbase:        coinbase,
		wrapToken:       wrapToken,
		slippage:        slippage,
		rpcUrl:          rpcUrl,
		privK:           privK,
		chainId:         chainId,
		coinMap:         coinMap,
		contractAddress: contractAddress,
	}, nil
}
