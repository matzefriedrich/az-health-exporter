package monitor

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/goccy/go-yaml"
)

type Config struct {
	Environment EnvConfig
	Resources   ResourceConfig
}

type EnvConfig struct {
	SubscriptionID             string `env:"AZURE_SUBSCRIPTION_ID,required"`
	PollInterval               int    `env:"POLL_INTERVAL_SECONDS"`
	ResourcesConfigurationFile string `env:"RESOURCES_CONFIG_FILE,required,file"`
	TenantID                   string `env:"AZURE_TENANT_ID,required"`
	ClientID                   string `env:"AZURE_CLIENT_ID,required"`
	ClientSecret               string `env:"AZURE_CLIENT_SECRET,required"`
}

type ResourceConfig struct {
	Resources []Resource `yaml:"resources"`
}

type Resource struct {
	ResourceGroup string `yaml:"resource_group"`
	Name          string `yaml:"name"`
	Type          string `yaml:"type"`
}

func LoadConfig() (*Config, error) {

	config := EnvConfig{}
	envOptions := env.Options{}
	configErr := env.ParseWithOptions(&config, envOptions)
	if configErr != nil {
		return nil, fmt.Errorf("failed to parse config: %w", configErr)
	}

	resourceConfig := ResourceConfig{}
	yamlBuffer := []byte(config.ResourcesConfigurationFile)
	yamlErr := yaml.Unmarshal(yamlBuffer, &resourceConfig)
	if yamlErr != nil {
		return nil, fmt.Errorf("failed to parse resources config: %w", yamlErr)
	}

	return &Config{
		Environment: config,
		Resources:   resourceConfig,
	}, nil
}
