package x

import (
	"fmt"

	"github.com/spf13/viper"
)

var searchCfgPath = []string{
	"../ecmauto",
	"C:/Program Files/Common Files/CAXA Shared/CAXA EAP/1.0/Bin/",
	"./alt_cfg/",
	"./config_client/",
	"../config-srv-api/config/",
	"../config-srv-api/",
	"./config-srv-api/",
	"./configd/",
	"./cfg/",
	"./web/",
	"$HOME",
	".",
}

// 获取 viper 配置,根据文件名(不含后缀)
func ViperCfg(cfgFileName string, cfgType string, funcSetDefault func(*viper.Viper)) (*viper.Viper, error) {
	vip := viper.New()
	if funcSetDefault != nil {
		funcSetDefault(vip)
	}
	for _, path := range searchCfgPath {
		vip.AddConfigPath(path)
	}
	vip.SetConfigName(cfgFileName)
	if cfgType != "" {
		vip.SetConfigType(cfgType) // json > toml > yaml > properties
	}
	if err := vip.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("$config file not found: %w", err)
		} else {
			return nil, fmt.Errorf("$config file was found but: %w", err)
		}
	}
	return vip, nil
}

// 获取 viper 配置,根据文件名(不含后缀)
func ViperSimpleCfg(cfgFileName string) (*viper.Viper, error) {
	vip := viper.New()

	for _, path := range searchCfgPath {
		vip.AddConfigPath(path)
	}
	vip.SetConfigName(cfgFileName)

	if err := vip.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("$config file not found: %w", err)
		} else {
			return nil, fmt.Errorf("$config file was found but: %w", err)
		}
	}
	return vip, nil
}
