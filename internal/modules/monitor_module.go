package modules

import (
	"log"

	"github.com/matzefriedrich/az-health-exporter/internal/monitor"
	"github.com/matzefriedrich/parsley/pkg/features"
	"github.com/matzefriedrich/parsley/pkg/types"
)

func MonitorModule(registry types.ServiceRegistry) error {
	features.RegisterLazy[monitor.HealthMonitor](registry, newHealthMonitorFactory, types.LifetimeSingleton)
	return nil
}

func newHealthMonitorFactory() monitor.HealthMonitor {
	config, configErr := monitor.LoadConfig()
	if configErr != nil {
		log.Fatalf("failed to load config: %v", configErr)
	}
	monitor, err := monitor.NewHealthMonitor(config)
	if err != nil {
		log.Fatalf("failed to create health monitor: %v", err)
	}
	return monitor
}
