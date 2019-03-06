package config

import (
	"errors"
	"strings"

	"github.com/spf13/viper"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/logger"
)

var Cfg *Config

const (
	LogLevelEnvKey  = "LOG_LEVEL"
	LogLevelDefault = logger.INFO
)

type Config struct {
	Node Node   `mapstructure:"NODE"`
	Db   DB     `mapstructure:"DB"`
	Port string `mapstructure:"PORT"`
}

// Load set configuration parameters.
// At first read config from file
// After that read environment variables
func Load(defaultConfigPath string) error {
	cfg, err := read(defaultConfigPath)
	if err != nil {
		return err
	}
	Cfg = cfg
	return validate()
}

func read(defaultConfigPath string) (*Config, error) {
	cfg := new(Config)

	// read config from file - it will be default values
	viper.SetConfigFile(defaultConfigPath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	// read parameters from environment variables -> they override default values from file
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func validate() error {
	if len(Cfg.Node.Host) == 0 {
		return errors.New("NODE_HOST parameter is empty")
	}
	if len(Cfg.Db.Host) == 0 {
		return errors.New("DB_HOST parameter is empty")
	}
	if len(Cfg.Db.Name) == 0 {
		return errors.New("DB_NAME parameter is empty")
	}
	if len(Cfg.Node.ChainType) == 0 {
		return errors.New("CHAIN parameter is empty")
	}
	return nil
}
