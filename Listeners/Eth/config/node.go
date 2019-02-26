package config

type Node struct {
	Host             string 	`mapstructure:"HOST"`
	StartBlockHeight string  	`mapstructure:"STARTBLOCK"`
	Confirmations    string  	`mapstructure:"CONFIRMATIONS"`
	Ticker			 string		`mapstructure:"TICKER"`
}
