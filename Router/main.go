package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/wavesplatform/GatewaysInfrastructure/Router/clientgrpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Router/config"
	"github.com/wavesplatform/GatewaysInfrastructure/Router/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Router/server"
	"github.com/wavesplatform/GatewaysInfrastructure/Router/service"
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
	log.Infof("router will be started with configuration %+v", *config.Cfg)
	ctx := context.Background()
	ctx = logger.ToContext(ctx, log)
	if err := clientgrpc.InitAllGrpcClients(ctx, config.Cfg); err != nil {
		log.Fatal("Can't init grpc clients: ", err)
	}

	blockchainService := service.New(
		clientgrpc.GetEthAdapterClient(),
		clientgrpc.GetWavesAdapterClient(),
		clientgrpc.GetEthListenerClient(),
		clientgrpc.GetWavesListenerClient())

	if err := server.InitAndStart(ctx, config.Cfg.GrpcPort, blockchainService); err != nil {
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
