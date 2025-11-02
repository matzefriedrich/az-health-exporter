package monitor

import (
	"context"
	"fmt"
	"iter"
	"log"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resourcehealth/armresourcehealth"
)

// healthMonitor manages resource health monitoring
type healthMonitor struct {
	client            *armresourcehealth.AvailabilityStatusesClient
	config            *Config
	healthStatus      map[string]*ResourceHealth
	mu                sync.RWMutex
	prometheusMetrics *PrometheusMetrics
}

type HealthMonitor interface {
	StartMonitoring(ctx context.Context)
	GetHealthStatus(_ context.Context) ([]*ResourceHealth, error)
}

// NewHealthMonitor creates a new health monitor instance
func NewHealthMonitor(config *Config) (HealthMonitor, error) {

	environment := config.Environment
	credential, err := azidentity.NewClientSecretCredential(environment.TenantID, environment.ClientID, environment.ClientSecret, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create credential: %w", err)
	}

	subscriptionID := environment.SubscriptionID
	client, err := armresourcehealth.NewAvailabilityStatusesClient(subscriptionID, credential, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	prometheusMetrics, _ := registerPrometheusMetrics()

	return &healthMonitor{
		client:            client,
		config:            config,
		healthStatus:      make(map[string]*ResourceHealth),
		prometheusMetrics: prometheusMetrics,
	}, nil
}

// StartMonitoring begins periodic health checks
func (m *healthMonitor) StartMonitoring(ctx context.Context) {

	interval := m.config.Environment.PollInterval
	if interval == 0 {
		interval = 60
	}
	duration := time.Duration(interval) * time.Second
	log.Printf("Starting health monitoring (interval: %v)", duration)

	notificationContext, cancel := context.WithCancel(ctx)
	defer cancel()

	m.checkAllResources(ctx)

	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for {
		select {
		case <-notificationContext.Done():
			log.Println("Monitoring stopped")
			return
		case <-ticker.C:
			m.checkAllResources(notificationContext)
		}
	}
}

func (m *healthMonitor) getConfiguredResources(ctx context.Context) iter.Seq[*ResourceInfo] {
	return func(yield func(*ResourceInfo) bool) {
		if m.config.Resources.Resources == nil {
			return
		}
		cancellationContext, cancel := context.WithCancel(ctx)
		defer cancel()
		for _, next := range m.config.Resources.Resources {
			select {
			case <-cancellationContext.Done():
				return
			default:
			}
			info := NewResourceInfo(m.config.Environment.SubscriptionID, next)
			yield(info)
		}
	}
}

func (m *healthMonitor) checkAllResources(ctx context.Context) {
	log.Println("Checking resource health...")
	resources := m.getConfiguredResources(ctx)
	for resource := range resources {
		health, err := m.checkResourceHealth(ctx, resource)
		if err != nil {
			log.Printf("Error checking %s: %v", resource.Name(), err)
			continue
		}
		m.updateHealthStatus(health)
	}
}

// checkResourceHealth retrieves health status for a single resource
func (m *healthMonitor) checkResourceHealth(ctx context.Context, resource *ResourceInfo) (*ResourceHealth, error) {

	id := resource.ID()
	resourceResponse, err := m.client.GetByResource(ctx, id, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get health status: %w", err)
	}

	status := resourceResponse.AvailabilityStatus
	availabilityState := "Unknown"
	if status.Properties.AvailabilityState != nil {
		availabilityState = string(*status.Properties.AvailabilityState)
	}

	healthy := availabilityState == "Available"

	summary := ""
	if status.Properties.Summary != nil {
		summary = *status.Properties.Summary
	}

	reasonType := ""
	if status.Properties.ReasonType != nil {
		reasonType = *status.Properties.ReasonType
	}

	return &ResourceHealth{
		ID:                id,
		Name:              resource.Name(),
		Type:              resource.Type(),
		AvailabilityState: availabilityState,
		Summary:           summary,
		ReasonType:        reasonType,
		ResourceGroup:     resource.ResourceGroup(),
		LastUpdated:       time.Now(),
		Healthy:           healthy,
	}, nil
}

func (m *healthMonitor) GetHealthStatus(_ context.Context) ([]*ResourceHealth, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	results := make([]*ResourceHealth, 0)
	for _, status := range m.healthStatus {
		results = append(results, status)
	}
	return results, nil
}

func (m *healthMonitor) updateHealthStatus(health *ResourceHealth) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.healthStatus[health.ID] = health

	// Update Prometheus metrics
	healthValue := 0.0
	if health.Healthy {
		healthValue = 1.0
	}

	m.prometheusMetrics.healthyGauge.WithLabelValues(
		withResourceGroup(health.ResourceGroup),
		withResourceID(health.ID),
		withResourceName(health.Name),
		withResourceType(health.Type),
		withAvailabilityState(health.AvailabilityState)).
		Set(healthValue)

	lastUpdateUnixTimestamp := float64(health.LastUpdated.Unix())
	m.prometheusMetrics.lastCheckGauge.WithLabelValues(
		withResourceID(health.ID),
		withResourceName(health.Name)).
		Set(lastUpdateUnixTimestamp)

	log.Printf("Updated: %s - %s (%s)",
		health.Name,
		health.AvailabilityState,
		health.Summary)
}
