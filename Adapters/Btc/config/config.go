package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/spf13/viper"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/logger"
)

const (
	LogLevelEnvKey  = "LOG_LEVEL"
	LogLevelDefault = logger.INFO
	testnetNetwork  = "testnet"
	mainnetNetwork  = "mainnet"
)

var (
	btcMainnetChainParams = chaincfg.MainNetParams
	btcTestnetChainParams = chaincfg.TestNet3Params
)

var Cfg *Config

type Config struct {
	Node        Node        `mapstructure:"NODE"`
	Port        string      `mapstructure:"PORT"`
	Decimals    int         `mapstructure:"DECIMALS"`
	Db          DB          `mapstructure:"DB"`
	DataService HttpService `mapstructure:"DATASERVICE"`
}

type HttpService struct {
	Url string `mapstructure:"URL"`
}

type Node struct {
	Host         string `mapstructure:"HOST"`
	User         string `mapstructure:"USER"`
	Password     string `mapstructure:"PASSWORD"`
	HTTPPostMode bool   `mapstructure:"HTTPPOSTMODE"`
	DisableTLS   bool   `mapstructure:"DISABLETLS"`
	ChainType    string `mapstructure:"CHAINTYPE"`
	FeeRateMax   uint64 `mapstructure:"FEE_RATE_MAX"`
	ChainParams  *chaincfg.Params
}

type DB struct {
	Name string `mapstructure:"NAME"`
	Host string `mapstructure:"HOST"`
}

func btcParamsByChainType(chainType string) *chaincfg.Params {
	switch chainType {
	case mainnetNetwork:
		return &btcMainnetChainParams
	case testnetNetwork:
		return &btcTestnetChainParams
	default:
		return nil
	}
}

func (c *Config) String() string {
	return fmt.Sprintf("NODE_HOST: %s, DATA_URL: %s, USER: %v, PORT: %s",
		c.Node.Host, c.DataService.Url, c.Node.User, c.Port,
	)
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
	if err := validateConfig(); err != nil {
		return err
	}
	cfg.Node.ChainParams = btcParamsByChainType(cfg.Node.ChainType)
	return nil
}

func validateConfig() error {
	if len(Cfg.Node.Host) == 0 {
		return errors.New("NODE_HOST parameter is empty")
	}
	if len(Cfg.Port) == 0 {
		return errors.New("PORT parameter is empty")
	}
	if len(Cfg.Node.ChainType) == 0 {
		return errors.New("CHAINTYPE parameter is empty")
	}
	return nil
}
