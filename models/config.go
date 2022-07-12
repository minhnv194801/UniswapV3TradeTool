package models

type Config struct {
	AppConfig     AppConfig         `json:"app_config"`
	DEXSetting    DEXSetting        `json:"dex_setting"`
	CoinContracts map[string]string `json:"coin_contracts"`
}

type AppConfig struct {
	RunExchange        bool `json:"run_exchange"`
	RunDex             bool `json:"run_dex"`
	RunTrustKeysCaller bool `json:"run_trust_keys_caller"`
}

type DEXSetting struct {
	BaseUrl         string  `json:"base_url"`
	ContractAddress string  `json:"contract_address"`
	QuoterAddress   string  `json:"quoter_address"`
	FactoryAddress  string  `json:"factory_address"`
	CoinbaseName    string  `json:"coinbase_name"`
	WrapToken       string  `json:"wrap_token"`
	Slippage        float64 `json:"slippage"`
	PrivateKey      string  `json:"private_key"`
	ChainId         int64   `json:"chain_id"`
	TimeStep        int64   `json:"time_step"`
}

type CoinsMap map[string]string
