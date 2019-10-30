package config

import (
	"fmt"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/models"
)

type Node struct {
	Host             string           `mapstructure:"HOST"`
	StartBlockHeight uint64           `mapstructure:"STARTBLOCK"`
	Confirmations    uint64           `mapstructure:"CONFIRMATIONS"`
	ChainType        models.ChainType `mapstructure:"CHAIN"`
	ApiKey           string           `mapstructure:"APIKEY"`
	RequestTimeoutMs int64            `mapstructure:"REQUEST_TIMEOUT_MS"`
}

func (n *Node) String() string {
	return fmt.Sprintf("NODE_HOST: %s, NODE_STARTBLOCK: %d, NODE_CONFIRMATIONS: %d, NODE_CHAIN: %s, NODE_REQUEST_TIMEOUT_MS: %d,",
		n.Host, n.StartBlockHeight, n.Confirmations, n.ChainType, n.RequestTimeoutMs)
}
