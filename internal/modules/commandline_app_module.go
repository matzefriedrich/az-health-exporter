package modules

import (
	"github.com/matzefriedrich/az-health-exporter/internal"
	"github.com/matzefriedrich/az-health-exporter/internal/commands"
	"github.com/matzefriedrich/az-health-exporter/internal/monitor"
	"github.com/matzefriedrich/cobra-extensions/pkg/charmer"
	"github.com/matzefriedrich/parsley/pkg/features"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/types"
)

func CommandlineAppModule(registry types.ServiceRegistry) error {
	_ = registration.RegisterSingleton(registry, configureCommandlineApplication)
	return nil
}

func configureCommandlineApplication(
	monitorFactory features.Lazy[monitor.HealthMonitor]) *charmer.CommandLineApplication {

	name := internal.GetInformativeApplicationName()
	app := charmer.NewCommandLineApplication(name, "")

	app.AddCommand(commands.NewHealthMonitorCommand(monitorFactory))

	return app
}
