package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	RPCWS     string `mapstructure:"rpc_ws"`
	Contract  string `mapstructure:"contract"`
	WorkerNum int    `mapstructure:"worker_num"`

	DB struct {
		DSN string `mapstructure:"dsn"`
	} `mapstructure:"db"`
}

func Load() Config {

	viper.SetConfigFile("configs/config.yaml")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	var c Config
	viper.Unmarshal(&c)

	return c
}
