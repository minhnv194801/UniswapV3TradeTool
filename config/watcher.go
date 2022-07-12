package config

import (
	"encoding/json"
	"orderbot/models"
	"orderbot/utils"
)

type WatcherConfig struct {
	dexConfig models.DEXSetting
	appConfig models.AppConfig
	coinsMap  models.CoinsMap
}

func (watcher *WatcherConfig) GetDexConf() models.DEXSetting {
	return watcher.dexConfig
}

func (watcher *WatcherConfig) GetCoinsMap() models.CoinsMap {
	return watcher.coinsMap
}

func NewWatcherConfig() (Watcher, error) {
	conf := &WatcherConfig{}
	err := conf.LoadConfig()
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func (watcher *WatcherConfig) LoadConfig() error {
	config, err := utils.ReadConfig()
	if err != nil {
		return err
	}
	b, err := json.Marshal(config)
	if err != nil {
		return err
	}
	// load app config
	conf := models.Config{}
	err = json.Unmarshal(b, &conf)
	if err != nil {
		return err
	}

	watcher.appConfig = conf.AppConfig
	watcher.dexConfig = conf.DEXSetting
	watcher.coinsMap = conf.CoinContracts

	return nil
}
