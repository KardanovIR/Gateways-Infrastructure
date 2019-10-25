package config

import "github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/models"

type Node struct {
	Host             string           `mapstructure:"HOST"`
	StartBlockHeight int64            `mapstructure:"STARTBLOCK"`
	Confirmations    string           `mapstructure:"CONFIRMATIONS"`
	ChainType        models.ChainType `mapstructure:"CHAIN"`
	ApiKey           string           `mapstructure:"APIKEY"`
}
