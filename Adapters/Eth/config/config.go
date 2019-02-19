package config

import (
	"errors"
	"github.com/spf13/viper"
	"strings"
)

var Cfg *Config

type Config struct {
	Node NodeConfig `mapstructure:"NODE"`
}

type NodeConfig struct {
	Host string `mapstructure:"HOST"`
}

// LoadConfig set configuration parameters.
// At first read config from file
// After that read environment variables
func LoadConfig(defaultConfigPath string) error {
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
	if len(Cfg.Node.Host) == 0 {
		return errors.New("NODE_HOST parameter is empty")
	}
	return nil
}
