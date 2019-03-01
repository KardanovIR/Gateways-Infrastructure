package config

type Node struct {
	Host             string `mapstructure:"HOST"`
	StartBlockHeight int64  `mapstructure:"STARTBLOCK"`
	Confirmations    string `mapstructure:"CONFIRMATIONS"`
	ChainType        string `mapstructure:"CHAIN"`
}
