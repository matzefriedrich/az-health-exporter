package main

import (
	"bufio"
	"context"
	"io/fs"
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
	_ = registry.RegisterModule(modules.CommandlineAppModule)
	_ = registry.RegisterModule(modules.MonitorModule)
	resolver := resolving.NewResolver(registry)

	ctx := context.Background()
	resolverContext := resolving.NewScopedContext(ctx)
	app, _ := resolving.ResolveRequiredService[*charmer.CommandLineApplication](resolverContext, resolver)

	err := app.Execute(resolverContext)
	if err != nil {
		log.Fatalf("Failed to execute command: %v", err)
	}
}

func printBanner() {
	bannerFile, _ := resources.Resources.Open(resources.BannerTxt)
	defer func(bannerFile fs.File) {
		_ = bannerFile.Close()
	}(bannerFile)
	scanner := bufio.NewScanner(bannerFile)
	for scanner.Scan() {
		line := scanner.Text()
		log.Println(line)
	}
}
