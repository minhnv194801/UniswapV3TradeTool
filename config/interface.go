package config

import "orderbot/models"

type Watcher interface {
	GetDexConf() models.DEXSetting
	GetCoinsMap() models.CoinsMap
}
