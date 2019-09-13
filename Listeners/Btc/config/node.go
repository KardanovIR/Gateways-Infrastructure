package config

import (
	"fmt"
)

type Node struct {
	Host             string           `mapstructure:"HOST"`
	StartBlockHeight uint64           `mapstructure:"STARTBLOCK"`
	Confirmations    uint64           `mapstructure:"CONFIRMATIONS"`
	User     	string `mapstructure:"USER"`
	Password      string             `mapstructure:"PASSWORD"`
	HTTPPostMode bool `mapstructure:"HTTPPOSTMODE"`
	DisableTLS bool `mapstructure:"DISABLETLS"`
}

func (n *Node) String() string {
	return fmt.Sprintf("NODE_HOST: %s, NODE_STARTBLOCK: %d, NODE_CONFIRMATIONS: %d, NODE_USER: %s",
		n.Host, n.StartBlockHeight, n.Confirmations, n.User)
}
