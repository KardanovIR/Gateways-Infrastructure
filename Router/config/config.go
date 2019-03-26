package config

import (
	"errors"
	"github.com/wavesplatform/GatewaysInfrastructure/Router/logger"
	"strings"

	"github.com/spf13/viper"
)

const (
	LogLevelEnvKey  = "LOG_LEVEL"
	LogLevelDefault = logger.INFO
)

var Cfg *Config

type Config struct {
	Adapters  BlockchainsList `mapstructure:"ADAPTER"`
	Listeners BlockchainsList `mapstructure:"LISTENER"`
	GrpcPort  string          `mapstructure:"PORT"`
}

type BlockchainsList struct {
	Eth   string `mapstructure:"ETH"`
	Waves string `mapstructure:"WAVES"`
}

// Load set configuration parameters.
// At first read config from file
// After that read environment variables
func Load(defaultConfigPath string) error {
	cfg := new(Config)

	// read config from file - it will be default values
	viper.SetConfigFile(defaultConfigPath)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(cfg); err != nil {
		return err
	}

	// read parameters from environment variables -> they override default values from file
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.Unmarshal(cfg); err != nil {
		return err
	}
	Cfg = cfg
	return validateConfig()
}

func validateConfig() error {
	if len(Cfg.Adapters.Waves) == 0 {
		return errors.New("ADAPTER_WAVES parameter is empty")
	}
	if len(Cfg.Listeners.Waves) == 0 {
		return errors.New("LISTENER_WAVES parameter is empty")
	}
	if len(Cfg.GrpcPort) == 0 {
		return errors.New("PORT parameter is empty")
	}
	return nil
}
