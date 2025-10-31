package main

import (
	"bufio"
	"context"
	"log"

	"github.com/matzefriedrich/az-health-exporter/internal"
	"github.com/matzefriedrich/az-health-exporter/internal/modules"
	"github.com/matzefriedrich/az-health-exporter/internal/resources"
	"github.com/matzefriedrich/cobra-extensions/pkg/charmer"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/resolving"
)

func main() {

	printBanner()
	log.Println(internal.GetInformativeApplicationName())

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

func printBanner() {
	bannerFile, _ := resources.Resources.Open(resources.BannerTxt)
	defer bannerFile.Close()
	scanner := bufio.NewScanner(bannerFile)
	for scanner.Scan() {
		line := scanner.Text()
		log.Println(line)
	}
}
