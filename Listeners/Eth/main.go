package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Eth/config"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Eth/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Eth/repositories"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Eth/services"
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
	log.Infof(fmt.Sprintf("Eth listener will be started with configuration %+v", config.Cfg))
	ctx := context.Background()
	ctx = logger.ToContext(ctx, log)

	repository, err := repositories.New(ctx, config.Cfg.Db.Host, config.Cfg.Db.Name)
	if err != nil {
		log.Fatal("Can't create db connection: ", err)
	}
	if err := services.New(ctx, config.Cfg.Node.Host, repository); err != nil {
		log.Fatal("Can't create node's client: ", err)
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