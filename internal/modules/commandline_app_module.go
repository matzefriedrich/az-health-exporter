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
	monitorFactory func() monitor.HealthMonitor) *charmer.CommandLineApplication {

	app := charmer.NewCommandLineApplication("", "")

	monitor := monitorFactory()
	app.AddCommand(commands.NewHealthMonitorCommand(monitor))

	return app
}
