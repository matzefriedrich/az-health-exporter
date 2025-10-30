package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"iter"
	"log"
	"net/http"
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
}

// NewHealthMonitor creates a new health monitor instance
func NewHealthMonitor(config *Config) (*healthMonitor, error) {

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create credential: %w", err)
	}

	subscriptionID := config.Environment.SubscriptionID
	client, err := armresourcehealth.NewAvailabilityStatusesClient(subscriptionID, cred, nil)
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

	environment := m.config.Environment
	duration := time.Duration(environment.PollInterval) * time.Second
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
	resp, err := m.client.GetByResource(ctx, id, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get health status: %w", err)
	}

	status := resp.AvailabilityStatus
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
		ResourceID:        id,
		ResourceName:      resource.Name(),
		ResourceType:      resource.Type(),
		AvailabilityState: availabilityState,
		Summary:           summary,
		ReasonType:        reasonType,
		LastUpdated:       time.Now(),
		Healthy:           healthy,
	}, nil
}

func (m *healthMonitor) updateHealthStatus(health *ResourceHealth) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.healthStatus[health.ResourceID] = health

	// Update Prometheus metrics
	healthValue := 0.0
	if health.Healthy {
		healthValue = 1.0
	}

	m.prometheusMetrics.healthyGauge.WithLabelValues(
		withResourceID(health.ResourceID),
		withResourceName(health.ResourceName),
		withResourceType(health.ResourceType),
		withAvailabilityState(health.AvailabilityState)).
		Set(healthValue)

	lastUpdateUnixTimestamp := float64(health.LastUpdated.Unix())
	m.prometheusMetrics.lastCheckGauge.WithLabelValues(
		withResourceID(health.ResourceID),
		withResourceName(health.ResourceName)).
		Set(lastUpdateUnixTimestamp)

	log.Printf("Updated: %s - %s (%s)",
		health.ResourceName,
		health.AvailabilityState,
		health.Summary)
}

// statusHandler returns current health status for all resources
func (m *healthMonitor) statusHandler(w http.ResponseWriter, r *http.Request) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	statuses := make([]*ResourceHealth, 0, len(m.healthStatus))
	for _, status := range m.healthStatus {
		statuses = append(statuses, status)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"timestamp": time.Now(),
		"resources": statuses,
	})
}
