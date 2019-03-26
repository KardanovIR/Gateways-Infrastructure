package clientgrpc

import (
	"context"
	"fmt"

	"github.com/wavesplatform/GatewaysInfrastructure/Router/config"
)

func InitAllGrpcClients(ctx context.Context, cfg *config.Config) error {
	if len(cfg.Adapters.Eth) > 0 {
		if err := NewEthAdapterClient(ctx, cfg.Adapters.Eth); err != nil {
			return fmt.Errorf("can't setup connection to eth adapter: %s", err)
		}
	}
	if len(cfg.Listeners.Eth) > 0 {
		if err := NewEthListenerClient(ctx, cfg.Listeners.Eth); err != nil {
			return fmt.Errorf("can't setup connection to eth listener: %s", err)
		}
	}
	if err := NewWavesListenerClient(ctx, cfg.Listeners.Waves); err != nil {
		return fmt.Errorf("can't setup connection to waves listener: %s", err)
	}
	if err := NewWavesAdapterClient(ctx, cfg.Adapters.Waves); err != nil {
		return fmt.Errorf("can't setup connection to waves adapter: %s", err)
	}
	return nil
}
