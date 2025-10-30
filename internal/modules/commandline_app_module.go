package modules

import (
	"github.com/matzefriedrich/az-health-exporter/internal/commands"
	"github.com/matzefriedrich/az-health-exporter/internal/monitor"
	"github.com/matzefriedrich/cobra-extensions/pkg/charmer"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/types"
)

func CommandlineAppModule(registry types.ServiceRegistry) error {
	registration.RegisterSingleton(registry, configureCommandlineApplication)
	return nil
}

func configureCommandlineApplication(
	monitor monitor.HealthMonitor) *charmer.CommandLineApplication {

	app := charmer.NewCommandLineApplication("az-health-exporter", "")

	app.AddCommand(commands.NewHealthMonitorCommand(monitor))

	return app
}
