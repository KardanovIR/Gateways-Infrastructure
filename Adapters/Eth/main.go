package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/GatewaysInfrastructure/Adapters/Eth/config"
	"github.com/GatewaysInfrastructure/Adapters/Eth/services"
)

func main() {
	var configPath string

	flag.StringVar(&configPath, "config-path", "./config/config.yml", "A path to config file")
	flag.Parse()

	if err := config.LoadConfig(configPath); err != nil {
		log.Fatal("Loading of configuration failed with error:", err)
	}
	log.Println(fmt.Sprintf("Eth adapter will be started with configuration %+v", config.Cfg))
	if err := services.New(config.Cfg.Node.Host); err != nil {
		log.Fatal("Can't create node's client", err)
	}
	ctx := context.Background()
	gas, err := services.GetNodeClient().SuggestGasPrice(ctx)
	fmt.Println("gas ", gas, "err", err)
}
