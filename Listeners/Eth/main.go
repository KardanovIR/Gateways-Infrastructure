package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/GatewaysInfrastructure/Listeners/Eth/config"
	"github.com/GatewaysInfrastructure/Listeners/Eth/repositories"
	"github.com/GatewaysInfrastructure/Listeners/Eth/services"
)

func main() {
	var configPath string

	flag.StringVar(&configPath, "config-path", "./config/config.yml", "A path to config file")
	flag.Parse()

	if err := config.LoadConfig(configPath); err != nil {
		log.Fatal("Loading of configuration failed with error:", err)
	}
	log.Println(fmt.Sprintf("Eth listener will be started with configuration %+v", config.Cfg))
	repository, err := repositories.New(config.Cfg.Db.Host, config.Cfg.Db.Name)
	if err != nil {
		log.Fatal("Can't create db connection: ", err)
	}
	if err := services.New(config.Cfg.Node.Host, repository); err != nil {
		log.Fatal("Can't create node's client: ", err)
	}
}
