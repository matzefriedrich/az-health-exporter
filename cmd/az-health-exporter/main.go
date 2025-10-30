package main

import (
	"context"
	"log"

	"github.com/matzefriedrich/az-health-exporter/internal/modules"
	"github.com/matzefriedrich/cobra-extensions/pkg/charmer"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/resolving"
)

func main() {

	registry := registration.NewServiceRegistry()
	registry.RegisterModule(modules.CommandlineAppModule)
	registry.RegisterModule(modules.MonitorModule)
	resolver := resolving.NewResolver(registry)

	ctx := context.Background()
	resolverContext := resolving.NewScopedContext(ctx)
	app, _ := resolving.ResolveRequiredService[*charmer.CommandLineApplication](resolver, resolverContext)

	err := app.Execute(resolverContext)
	if err != nil {
		log.Fatalf("Failed to execute command: %v", err)
	}

}
