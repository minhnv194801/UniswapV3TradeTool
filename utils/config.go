package utils

import "github.com/BurntSushi/toml"

func ReadConfig() (map[string]interface{}, error) {
	serviceConfig := map[string]interface{}{}

	// For testing need to adjust path to ../conf/app.conf
	_, err := toml.DecodeFile("./conf/app.conf", &serviceConfig)
	if err != nil {
		return map[string]interface{}{}, err
	}

	return serviceConfig, nil
}
