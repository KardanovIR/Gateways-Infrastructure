package config

import (
	"fmt"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/models"
)

type Node struct {
	Host             string           `mapstructure:"HOST"`
	StartBlockHeight uint64           `mapstructure:"STARTBLOCK"`
	Confirmations    uint64           `mapstructure:"CONFIRMATIONS"`
	User     	string `mapstructure:"USER"`
	ChainType        models.ChainType `mapstructure:"CHAIN"`
	Password      string             `mapstructure:"PASSWORD"`
	HTTPPostMode bool `mapstructure:"HTTPPOSTMODE"`
	DisableTLS bool `mapstructure:"DISABLETLS"`
}

func (n *Node) String() string {
	return fmt.Sprintf("NODE_HOST: %s, NODE_STARTBLOCK: %d, NODE_CONFIRMATIONS: %d, NODE_USER: %s",
		n.Host, n.StartBlockHeight, n.Confirmations, n.User)
}
