This bot get price and create order base on that price.  

Setting: (Exchange)
-

- EpochLength: seconds exp: 10s
- EpochNumber: number of epoch exp: 4
- StartRate: start rate base on base price exp:100%
- StepRate: the change between each order exp:1%
- StepNumber: number of order in each epoch exp: 5
- Total price = (epochNumber + 1) * moneyInEachEpoch  
  plus 1 because we need 1 epochSlot of money to place order
  

Setting: (DEX) (for Uniswap)
- 

- Setting in file: `/conf/app.conf`
- Dex setting:  
  base_url: RPC URL of node exp: "https://ropsten.infura.io/v3/9b11910966d3430e9846e504d5847593" (string)  
  contract_address: Swap Router contract address exp: "0xE592427A0AEce92De3Edee1F18E0157C05861564" (string)    
  quoter_address: Quoter contract address exp: "0xb27308f9F90D607463bb33eA1BeBb41C27CE5AB6" (string)
  factory_address: Factory contract address exp: "0x1F98431c8aD98523631AE4a59f267346ea31F984" (string)
  coinbase_name: Name of coinbase exp: "ETH" (string)  
  wrap_token: Name of coinbase wrap token exp: "WETH" (string)    
  slippage: slippage tolerance in percent exp: 0.5 (float)  
  private_key: your private key  (string)  
  chain_id: chain id exp: 1 (int)  
  time_step: time step length (second) (int) exp: 30
- Coin contracts:  
  List of coin name match address (string)
  
Usage:
-

- Get price:  
  Return the price of coin buy compare to coin sell

  `orderbot price --sell-coin={coin sell} --buy-coin={coin buy}`  
  example:  
  `orderbot price --sell-coin=BNB --buy-coin=USDT`

- GetBalance:  
  Return the balance of user with a specific coin  
  `orderbot balance --coin={coin}`  
  example:  
  `orderbot balance --coin=BNB`  

- CreateOrder:  
  Create an order  
  `orderbot order --sell-coin={coin sell} --buy-coin={coin buy} --spend={spend coin}`  
  example:  
  `orderbot order --sell-coin=BNB --buy-coin=USDT --spend=0.1`  

- Start:  
  Start the auto order bot, each times, bot spend an amount of `spend-each` coin to buy  
  Start bot order
  `orderbot start --sell-coin={coin sell} --buy-coin={coin buy} --stop-balance={balance remain that make stop immediately} --spend-each={coin spend each order} --min-price={min price}`  
  example:
  `orderbot start --sell-coin=BNB --buy-coin=USDT --stop-balance=4.1 --spend-each=0.01 --min-price=400` 
  
Warning:
-

- Only support ERC20 and coinbase
- Don't set the time step too short if you sell ERC20 (> 15s)
- Uniswap trade tool isn't optimized