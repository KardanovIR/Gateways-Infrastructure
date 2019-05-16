package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/config"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/repositories"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/server"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/services"
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
		log.Fatal("Loading of configuration failed with error:", err)
	}
	log.Infof(fmt.Sprintf("Waves listener will be started with configuration %+v", config.Cfg))
	ctx := context.Background()
	ctx = logger.ToContext(ctx, log)

	err = repositories.New(ctx, config.Cfg.Db.Host, config.Cfg.Db.Name)
	if err != nil {
		log.Fatal("Can't create db connection: ", err)
	}

	if err := services.NewCallbackService(ctx, config.Cfg.CallbackUrl, config.Cfg.Node.ChainType); err != nil {
		log.Fatal("Can't create callback service: ", err)
	}

	repository := repositories.GetRepository()
	if err := services.New(ctx, &config.Cfg.Node, repository); err != nil {
		log.Fatal("Can't create node's client: ", err)
	}
	nodeReader := services.GetNodeReader()

	err = nodeReader.Start(ctx)
	if err != nil {
		log.Fatal("Can't start node reader: ", err)
	}
	defer nodeReader.Stop(ctx)

	if err := server.InitAndStart(ctx, config.Cfg.Port, repository); err != nil {
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
