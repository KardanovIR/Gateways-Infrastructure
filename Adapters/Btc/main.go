package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/repositories"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/services/data_client"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/services/node_client"
	"os"

	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/config"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/server"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config-path", "./config/config.yml", "A path to config file")
	isDebugMode := flag.Bool("debug", false, "debug mode")
	flag.Parse()

	log, err := initLogger(*isDebugMode)
	if err != nil {
		fmt.Println("Can't initialize logger", err)
		return
	}

	if err := config.Load(configPath); err != nil {
		log.Fatal("loading of configuration failed with error:", err)
	}
	log.Infof("btc adapter will be started with configuration %s", config.Cfg.String())
	ctx := context.Background()
	ctx = logger.ToContext(ctx, log)

	if err := repositories.New(ctx, config.Cfg.Db); err != nil {
		log.Fatal("can't create repository: ", err)
	}
	if err := node_client.New(ctx, config.Cfg.Node, repositories.GetRepository()); err != nil {
		log.Fatal("can't create node's client: ", err)
	}

	if err := data_client.NewDataClient(ctx, config.Cfg.DataService); err != nil {
		log.Fatal("can't create data's client: ", err)
	}

	log.Info("")
	if err := server.InitAndStart(ctx, config.Cfg.Port, node_client.GetNodeClient(), data_client.GetDataClient()); err != nil {
		log.Fatal("Can't start grpc server", err)
	}
}

// initLogger initializes logger: create logger, set logger format: json or text.
// text is used if application was started with flag '-debug'
// set log level according to environment variable LOG_LEVEL,
// if LOG_LEVEL was not set it uses INFO by default,
// if application was started with flag '-debug' set DEBUG level
func initLogger(isDebug bool) (logger.ILogger, error) {
	var level = config.LogLevelDefault
	if isDebug {
		level = logger.DEBUG
	}

	if l, ok := os.LookupEnv(config.LogLevelEnvKey); ok {
		level = logger.Level(l)
	}
	return logger.Init(!isDebug, level)
}
