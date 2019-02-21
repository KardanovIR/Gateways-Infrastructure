package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/config"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/server"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/services"
	"log"
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

	if err := server.New(config.Cfg.Port, services.GetNodeClient()); err != nil {
		log.Fatal("Can't create grpc server", err)
	}
	if err := server.GetGrpsServer().Start(); err != nil {
		log.Fatal("Can't start grpc server", err)
	}
}
