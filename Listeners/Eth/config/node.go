package config

type Node struct {
	Host             string `mapstructure:"HOST"`
	StartBlockHeight int64  `mapstructure:"STARTBLOCK"`
	Confirmations    int64  `mapstructure:"CONFIRMATIONS"`
}
